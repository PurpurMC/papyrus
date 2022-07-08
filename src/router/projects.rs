use actix_web::{get, HttpResponse};
use actix_web::web::{Data, Path, ServiceConfig};
use serde_json::json;
use sqlx::SqlitePool;
use crate::models::{Project, Version};
use crate::utils;

pub fn routes(config: &mut ServiceConfig) {
    config.service(list_projects);
    config.service(get_project);
}

#[get("")]
async fn list_projects(pool: Data<SqlitePool>) -> HttpResponse {
    let mut projects = match Project::all(&pool).await {
        Ok(projects) => projects,
        Err(err) => return HttpResponse::InternalServerError().json(json!({ "error": err.to_string() })),
    };

    projects.sort_by(|a, b| a.created_at.cmp(&b.created_at));
    let projects = projects.iter().map(|project| &project.name).collect::<Vec<_>>();

    HttpResponse::Ok().json(json!({
        "projects": projects,
    }))
}

#[get("/{project}")]
async fn get_project(pool: Data<SqlitePool>, project: Path<String>) -> HttpResponse {
    let project = match utils::router::project(&pool, &project.into_inner()).await {
        Ok(project) => project,
        Err(err) => return err,
    };

    let mut versions = match Version::all(&pool, &project.id).await {
        Ok(versions) => versions,
        Err(err) => return HttpResponse::InternalServerError().json(json!({ "error": err.to_string() })),
    };

    versions.sort_by(|a, b| a.created_at.cmp(&b.created_at));
    let versions = versions.iter().map(|versions| &versions.name).collect::<Vec<&String>>();

    HttpResponse::Ok().json(json!({
        "project": project.name,
        "versions": versions,
    }))
}
