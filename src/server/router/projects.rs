use crate::types::models::{Project, Version};
use crate::types::response::{ProjectResponse, ProjectsResponse};
use crate::types::Result;
use crate::utils::{response, verify};
use crate::SqlitePool;
use actix_web::web::{Data, Path, ServiceConfig};
use actix_web::{get, HttpResponse};

pub fn routes(config: &mut ServiceConfig) {
    config.service(all_projects);
    config.service(get_project);
}

#[get("")]
pub async fn all_projects(pool: Data<SqlitePool>) -> Result<HttpResponse> {
    let projects = Project::all(&pool).await?;
    let projects: Vec<String> = projects
        .iter()
        .map(|project| project.name.clone())
        .collect();

    response(ProjectsResponse { projects })
}

#[get("/{project}")]
pub async fn get_project(pool: Data<SqlitePool>, path: Path<String>) -> Result<HttpResponse> {
    let project = path.into_inner();
    let project = verify(Project::find_one(&project, &pool).await?)?;

    let mut versions = Version::find_all(&project.id, &pool).await?;
    versions.sort_by(|a, b| a.created_at.cmp(&b.created_at));
    let versions: Vec<String> = versions
        .iter()
        .map(|version| version.name.clone())
        .collect();

    response(ProjectResponse {
        project: project.name,
        versions,
    })
}
