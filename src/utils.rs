use std::io;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum Error {
    #[error("Failed to parse config file: {0}")]
    ConfigParseError(#[from] serde_json::Error),

    #[error("Failed to access config file: {0}")]
    ConfigFileError(#[from] io::Error),

    #[error("Failed to load config: {0}")]
    DatabaseError(#[from] sqlx::Error),
}

pub mod router {
    use crate::models::{Build, Project, Version};
    use actix_web::HttpResponse;
    use serde_json::json;
    use sqlx::SqlitePool;

    pub async fn project(pool: &SqlitePool, name: &String) -> Result<Project, HttpResponse> {
        let project = match Project::get(&pool, name).await {
            Ok(project) => project,
            Err(err) => {
                return Err(
                    HttpResponse::InternalServerError().json(json!({ "error": err.to_string() }))
                )
            }
        };

        match project {
            Some(project) => Ok(project),
            None => Err(HttpResponse::NotFound().json(json!({ "error": "Project not found" }))),
        }
    }

    pub async fn version(
        pool: &SqlitePool,
        project: &Project,
        version: &String,
    ) -> Result<Version, HttpResponse> {
        let version = match Version::get(&pool, version, &project.id).await {
            Ok(version) => version,
            Err(err) => {
                return Err(
                    HttpResponse::InternalServerError().json(json!({ "error": err.to_string() }))
                )
            }
        };

        match version {
            Some(version) => Ok(version),
            None => Err(HttpResponse::NotFound().json(json!({ "error": "Version not found" }))),
        }
    }

    pub async fn build(
        pool: &SqlitePool,
        version: &Version,
        build: &String,
    ) -> Result<Build, HttpResponse> {
        if build == "latest" {
            let builds = match Build::all(&pool, &version.id).await {
                Ok(builds) => builds,
                Err(err) => {
                    return Err(HttpResponse::InternalServerError()
                        .json(json!({ "error": err.to_string() })))
                }
            };

            match builds
                .iter()
                .filter(|build| build.hash.is_some())
                .rev()
                .find(|build| build.result != "FAILURE")
            {
                Some(build) => Ok(build.clone()),
                None => Err(HttpResponse::NotFound().json(json!({ "error": "Build not found" }))),
            }
        } else {
            let build = match Build::get(&pool, build, &version.id).await {
                Ok(build) => build,
                Err(err) => {
                    return Err(HttpResponse::InternalServerError()
                        .json(json!({ "error": err.to_string() })))
                }
            };

            match build {
                Some(build) => Ok(build),
                None => Err(HttpResponse::NotFound().json(json!({ "error": "Build not found" }))),
            }
        }
    }
}
