use crate::utils::Error;
use chrono::NaiveDateTime;
use serde::{Deserialize, Serialize};
use sqlx::SqlitePool;

#[derive(Serialize, Deserialize)]
pub struct Project {
    pub id: String,
    pub name: String,
    pub created_at: NaiveDateTime,
}

impl Project {
    pub async fn all(pool: &SqlitePool) -> Result<Vec<Project>, Error> {
        Ok(sqlx::query_as!(Project, r#"SELECT * FROM projects"#)
            .fetch_all(pool)
            .await?)
    }

    pub async fn get(pool: &SqlitePool, name: &String) -> Result<Option<Project>, Error> {
        Ok(
            sqlx::query_as!(Project, r#"SELECT * FROM projects WHERE name = ?"#, name)
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
    pub async fn all(pool: &SqlitePool, project: &String) -> Result<Vec<Version>, Error> {
        Ok(sqlx::query_as!(
            Version,
            r#"SELECT * FROM versions WHERE project_id = ?"#,
            project
        )
        .fetch_all(pool)
        .await?)
    }

    pub async fn get(
        pool: &SqlitePool,
        name: &String,
        project: &String,
    ) -> Result<Option<Version>, Error> {
        Ok(sqlx::query_as!(
            Version,
            r#"SELECT * FROM versions WHERE name = ? AND project_id = ?"#,
            name,
            project
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
    pub hash: Option<String>,
    pub file_extension: Option<String>,
    pub created_at: NaiveDateTime,
}

impl Build {
    pub async fn all(pool: &SqlitePool, version: &String) -> Result<Vec<Build>, Error> {
        Ok(sqlx::query_as!(
            Build,
            r#"SELECT * FROM builds WHERE version_id = ?"#,
            version
        )
        .fetch_all(pool)
        .await?)
    }

    pub async fn get(
        pool: &SqlitePool,
        name: &String,
        version: &String,
    ) -> Result<Option<Build>, Error> {
        Ok(sqlx::query_as!(
            Build,
            r#"SELECT * FROM builds WHERE name = ? AND version_id = ?"#,
            name,
            version
        )
        .fetch_optional(pool)
        .await?)
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
    pub async fn get(pool: &SqlitePool, build: &String) -> Result<Vec<Commit>, Error> {
        Ok(
            sqlx::query_as!(Commit, r#"SELECT * FROM commits WHERE build_id = ?"#, build)
                .fetch_all(pool)
                .await?,
        )
    }
}
