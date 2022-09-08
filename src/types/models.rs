use crate::types::request::{CreateBuildRequest, CreateBuildRequestCommit};
use crate::types::response::{BuildResponse, BuildResponseCommit};
use crate::types::{Error, Result};
use crate::SqlitePool;
use chrono::NaiveDateTime;
use nanoid::nanoid;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct Project {
    pub id: String,
    pub name: String,
    pub created_at: NaiveDateTime,
}

impl Project {
    pub async fn create(name: &String, pool: &SqlitePool) -> Result<Self> {
        let id = nanoid!();
        sqlx::query!(
            "INSERT OR IGNORE INTO projects (id, name) VALUES (?, ?)",
            id,
            name
        )
        .execute(pool)
        .await?;

        match Project::find_one(&name, pool).await? {
            Some(project) => Ok(project),
            None => {
                Project::delete(&name, pool).await?;
                Err(Error::InternalServerError)
            }
        }
    }

    pub async fn delete(name: &String, pool: &SqlitePool) -> Result<()> {
        sqlx::query!("DELETE FROM projects WHERE name = ?", name)
            .execute(pool)
            .await?;
        Ok(())
    }

    pub async fn all(pool: &SqlitePool) -> Result<Vec<Self>> {
        Ok(sqlx::query_as!(Project, "SELECT * FROM projects")
            .fetch_all(pool)
            .await?)
    }

    pub async fn find_one(name: &String, pool: &SqlitePool) -> Result<Option<Self>> {
        Ok(
            sqlx::query_as!(Project, "SELECT * FROM projects WHERE name = ?", name)
                .fetch_optional(pool)
                .await?,
        )
    }
}

#[derive(Serialize, Deserialize)]
pub struct Version {
    pub id: String,
    pub name: String,
    pub project_id: String,
    pub created_at: NaiveDateTime,
}

impl Version {
    pub async fn create(project_id: &String, name: &String, pool: &SqlitePool) -> Result<Self> {
        let id = nanoid!();
        sqlx::query!(
            "INSERT OR IGNORE INTO versions (id, project_id, name) VALUES (?, ?, ?)",
            id,
            project_id,
            name
        )
        .execute(pool)
        .await?;

        match Version::find_one(&project_id, &name, pool).await? {
            Some(version) => Ok(version),
            None => {
                Version::delete(&project_id, &name, pool).await?;
                Err(Error::InternalServerError)
            }
        }
    }

    pub async fn delete(project_id: &String, name: &String, pool: &SqlitePool) -> Result<()> {
        sqlx::query!(
            "DELETE FROM versions WHERE project_id = ? AND name = ?",
            project_id,
            name
        )
        .execute(pool)
        .await?;
        Ok(())
    }

    pub async fn find_all(project_id: &String, pool: &SqlitePool) -> Result<Vec<Self>>
    where
        Self: Sized,
    {
        Ok(sqlx::query_as!(
            Version,
            "SELECT * FROM versions WHERE project_id = ?",
            project_id
        )
        .fetch_all(pool)
        .await?)
    }

    pub async fn find_one(
        project_id: &String,
        name: &String,
        pool: &SqlitePool,
    ) -> Result<Option<Self>> {
        Ok(sqlx::query_as!(
            Version,
            "SELECT * FROM versions WHERE project_id = ? AND name = ?",
            project_id,
            name
        )
        .fetch_optional(pool)
        .await?)
    }
}

#[derive(Serialize, Deserialize, Clone)]
pub struct Build {
    pub id: String,
    pub name: String,
    pub version_id: String,
    pub result: String,
    pub duration: i64,
    pub timestamp: i64,
    pub created_at: NaiveDateTime,
}

impl Build {
    pub async fn create(
        version_id: &String,
        request: &CreateBuildRequest,
        pool: &SqlitePool,
    ) -> Result<Self> {
        let id = nanoid!();
        sqlx::query!("INSERT OR IGNORE INTO builds (id, version_id, name, result, duration, timestamp) VALUES (?, ?, ?, ?, ?, ?)", id, version_id, request.build, request.result, request.duration, request.timestamp).execute(pool).await?;

        match Build::find_one(&version_id, &request.build, pool).await? {
            Some(build) => Ok(build),
            None => {
                Build::delete(&version_id, &request.build, pool).await?;
                Err(Error::InternalServerError)
            }
        }
    }

    pub async fn delete(version_id: &String, name: &String, pool: &SqlitePool) -> Result<()> {
        sqlx::query!(
            "DELETE FROM builds WHERE version_id = ? AND name = ?",
            version_id,
            name
        )
        .execute(pool)
        .await?;
        Ok(())
    }

    pub async fn find_all(version_id: &String, pool: &SqlitePool) -> Result<Vec<Self>>
    where
        Self: Sized,
    {
        Ok(sqlx::query_as!(
            Build,
            "SELECT * FROM builds WHERE version_id = ?",
            version_id
        )
        .fetch_all(pool)
        .await?)
    }

    pub async fn find_one(
        version_id: &String,
        name: &String,
        pool: &SqlitePool,
    ) -> Result<Option<Self>> {
        Ok(sqlx::query_as!(
            Build,
            "SELECT * FROM builds WHERE version_id = ? AND name = ?",
            version_id,
            name
        )
        .fetch_optional(pool)
        .await?)
    }

    pub async fn to_response(
        &self,
        project: &String,
        version: &String,
        pool: &SqlitePool,
    ) -> Result<Option<BuildResponse>> {
        let commits = Commit::find_all(&self.id, pool).await?;
        let commits = commits
            .iter()
            .map(|commit| BuildResponseCommit {
                author: commit.author.clone(),
                email: commit.email.clone(),
                description: commit.description.clone(),
                hash: commit.hash.clone(),
                timestamp: commit.timestamp,
            })
            .collect();

        let hash = if &self.result == "SUCCESS" {
            match File::find_one(&self.id, pool).await? {
                Some(hash) => hash.hash,
                None => return Ok(None),
            }
        } else {
            "".into()
        };

        Ok(Some(BuildResponse {
            project: project.clone(),
            version: version.clone(),
            build: self.name.clone(),
            commits,
            result: self.result.clone(),
            md5: hash,
            duration: self.duration,
            timestamp: self.timestamp,
        }))
    }
}

#[derive(Serialize, Deserialize)]
pub struct Commit {
    pub id: String,
    pub build_id: String,
    pub author: String,
    pub email: String,
    pub description: String,
    pub hash: String,
    pub timestamp: i64,
    pub created_at: NaiveDateTime,
}

impl Commit {
    pub async fn create(
        build_id: &String,
        request: &CreateBuildRequestCommit,
        pool: &SqlitePool,
    ) -> Result<Self> {
        let id = nanoid!();
        sqlx::query!(
            "INSERT OR IGNORE INTO commits (id, build_id, author, email, description, hash, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?)",
            id,
            build_id,
            request.author,
            request.email,
            request.description,
            request.hash,
            request.timestamp
        )
            .execute(pool)
            .await?;

        match sqlx::query_as!(Commit, "SELECT * FROM commits WHERE id = ?", id)
            .fetch_optional(pool)
            .await?
        {
            Some(commit) => Ok(commit),
            None => {
                Commit::delete(&id, pool).await?;
                Err(Error::InternalServerError)
            }
        }
    }

    pub async fn delete(id: &String, pool: &SqlitePool) -> Result<()> {
        sqlx::query!("DELETE FROM commits WHERE id = ?", id)
            .execute(pool)
            .await?;
        Ok(())
    }

    pub async fn find_all(build_id: &String, pool: &SqlitePool) -> Result<Vec<Self>>
    where
        Self: Sized,
    {
        Ok(
            sqlx::query_as!(Commit, "SELECT * FROM commits WHERE build_id = ?", build_id)
                .fetch_all(pool)
                .await?,
        )
    }
}

#[derive(Serialize, Deserialize)]
pub struct File {
    pub id: String,
    pub build_id: String,
    pub hash: String,
    pub extension: String,
    pub created_at: NaiveDateTime,
}

impl File {
    pub async fn create(
        build_id: &String,
        extension: &String,
        hash: &String,
        pool: &SqlitePool,
    ) -> Result<Self> {
        let id = nanoid!();
        sqlx::query!(
            "INSERT OR IGNORE INTO files (id, build_id, hash, extension) VALUES (?, ?, ?, ?)",
            id,
            build_id,
            hash,
            extension
        )
        .execute(pool)
        .await?;

        match sqlx::query_as!(File, "SELECT * FROM files WHERE id = ?", id)
            .fetch_optional(pool)
            .await?
        {
            Some(file) => Ok(file),
            None => {
                File::delete(&id, pool).await?;
                Err(Error::InternalServerError)
            }
        }
    }

    pub async fn delete(id: &String, pool: &SqlitePool) -> Result<()> {
        sqlx::query!("DELETE FROM files WHERE id = ?", id)
            .execute(pool)
            .await?;
        Ok(())
    }

    pub async fn find_one(build_id: &String, pool: &SqlitePool) -> Result<Option<Self>>
    where
        Self: Sized,
    {
        Ok(
            sqlx::query_as!(File, "SELECT * FROM files WHERE build_id = ?", build_id)
                .fetch_optional(pool)
                .await?,
        )
    }
}

#[derive(Serialize, Deserialize)]
pub struct TempFile {
    pub id: String,
    pub extension: String,
    pub created_at: NaiveDateTime,
}

impl TempFile {
    pub async fn create(extension: &Option<String>, pool: &SqlitePool) -> Result<Self> {
        let id = nanoid!();
        let extension = extension.clone().unwrap_or("".into());

        sqlx::query!(
            "INSERT OR IGNORE INTO temp_files (id, extension) VALUES (?, ?)",
            id,
            extension
        )
        .execute(pool)
        .await?;
        match sqlx::query_as!(TempFile, "SELECT * FROM temp_files WHERE id = ?", id)
            .fetch_optional(pool)
            .await?
        {
            Some(temp_file) => Ok(temp_file),
            None => {
                TempFile::delete(&id, pool).await?;
                Err(Error::InternalServerError)
            }
        }
    }

    pub async fn delete(id: &String, pool: &SqlitePool) -> Result<()> {
        sqlx::query!("DELETE FROM temp_files WHERE id = ?", id)
            .execute(pool)
            .await?;
        Ok(())
    }

    pub async fn find_one(id: &String, pool: &SqlitePool) -> Result<Option<Self>>
    where
        Self: Sized,
    {
        Ok(
            sqlx::query_as!(TempFile, "SELECT * FROM temp_files WHERE id = ?", id)
                .fetch_optional(pool)
                .await?,
        )
    }
}
