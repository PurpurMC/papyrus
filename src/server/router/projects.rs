use crate::types::response::{ProjectResponse, ProjectsResponse};
use crate::types::Response;
use actix_web::get;
use actix_web::web::{Path, ServiceConfig};

pub fn routes(config: &mut ServiceConfig) {
    config.service(all_projects);
    config.service(get_project);
}

#[get("")]
pub async fn all_projects() -> Response<ProjectsResponse> {
    todo!()
}

#[get("/{project}")]
pub async fn get_project(path: Path<String>) -> Response<ProjectResponse> {
    let project = path.into_inner();

    todo!()
}
