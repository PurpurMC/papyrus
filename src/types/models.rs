use crate::types::response::{BuildResponse, BuildResponseCommit};
use crate::types::Result;
use crate::utils::verify;
use crate::SqlitePool;
use chrono::NaiveDateTime;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct Project {
    pub id: String,
    pub name: String,
    pub created_at: NaiveDateTime,
}

impl Project {
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
    ) -> Result<BuildResponse> {
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

        let file = verify(File::find_one(&self.id, pool).await?)?;

        Ok(BuildResponse {
            project: project.clone(),
            version: version.clone(),
            build: self.name.clone(),
            commits,
            result: self.result.clone(),
            md5: file.hash.clone(),
            duration: self.duration,
            timestamp: self.timestamp,
        })
    }
}

#[derive(Serialize, Deserialize)]
pub struct File {
    pub id: String,
    pub build_id: Option<String>,
    pub hash: String,
    pub extension: Option<String>,
    pub created_at: NaiveDateTime,
}

impl File {
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
pub struct Commit {
    pub id: String,
    pub build_id: String,
    pub author: String,
    pub email: String,
    pub description: String,
    pub hash: String,
    pub timestamp: i64,
}

impl Commit {
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
