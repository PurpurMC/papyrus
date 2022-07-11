use actix_web::{get, HttpResponse};
use actix_web::web::{Data, Path, ServiceConfig};
use serde_json::json;
use sqlx::SqlitePool;
use crate::models::{Build, Project, Version};
use crate::utils;

pub fn routes(config: &mut ServiceConfig) {
    config.service(get_version);
    config.service(get_version_detailed);
}

#[get("/{project}/{version}")]
async fn get_version(pool: Data<SqlitePool>, path: Path<(String, String)>) -> HttpResponse {
    let (project, version) = path.into_inner();
    let (project, version, mut builds) = match version_info(&pool, project, version).await {
        Ok((project, version, builds)) => (project, version, builds),
        Err(err) => return err,
    };

    builds.sort_by(|a, b| a.created_at.cmp(&b.created_at));

    let latest_build = builds.iter().filter(|build| build.hash.is_some()).rev().find(|build| build.result != "FAILURE").map(|build| &build.name);
    let builds = builds.iter().filter(|build| build.hash.is_some()).map(|build| &build.name).collect::<Vec<&String>>();

    HttpResponse::Ok().json(json!({
        "project": project.name,
        "version": version.name,
        "builds": {
            "latest": latest_build,
            "all": builds,
        }
    }))
}

#[get("/{project}/{version}/detailed")]
async fn get_version_detailed(pool: Data<SqlitePool>, path: Path<(String, String)>) -> HttpResponse {
    let (project, version) = path.into_inner();
    let (project, version, mut builds) = match version_info(&pool, project, version).await {
        Ok((project, version, builds)) => (project, version, builds),
        Err(err) => return err,
    };

    builds.sort_by(|a, b| a.created_at.cmp(&b.created_at));
    let mut res_builds = Vec::<crate::router::builds::Build>::new();

    for build in builds {
        if build.hash.is_none() {
            continue;
        }

        res_builds.push(crate::router::builds::Build {
            project: project.name.clone(),
            version: version.name.clone(),
            build: build.name.clone(),
            result: build.result.clone(),
            commits: match crate::router::builds::commits(&pool, &build).await {
                Ok(commits) => commits,
                Err(err) => return err,
            },
            md5: build.hash.unwrap_or("".into()),
            duration: build.duration,
            timestamp: build.timestamp,
        });
    }

    HttpResponse::Ok().json(json!({
        "project": project.name,
        "version": version.name,
        "builds": {
            "latest": res_builds.iter().rev().find(|build| build.result != "FAILURE"),
            "all": res_builds,
        }
    }))
}

async fn version_info(pool: &SqlitePool, project: String, version: String) -> Result<(Project, Version, Vec<Build>), HttpResponse> {
    let project = match utils::router::project(&pool, &project).await {
        Ok(projects) => projects,
        Err(err) => return Err(err),
    };

    let version = match utils::router::version(&pool, &project, &version).await {
        Ok(version) => version,
        Err(err) => return Err(err),
    };

    match Build::all(&pool, &version.id).await {
        Ok(builds) => Ok((project, version, builds)),
        Err(err) => Err(HttpResponse::InternalServerError().json(json!({ "error": err.to_string() })))
    }
}
