use crate::server::router::versions::{get_version, get_version_detailed};
use crate::types::response::{VersionResponse, VersionResponseDetailed};
use crate::types::Result;
use actix_web::web::{Path, ServiceConfig};
use actix_web::{get, HttpResponse};

pub fn routes(config: &mut ServiceConfig) {
    config.service(get_version);
    config.service(get_version_detailed);
}

#[get("/{project}/{version}/{build}")]
pub async fn get_build(path: Path<(String, String, String)>) -> Result<HttpResponse> {
    let (project, version, build) = path.into_inner();

    todo!()
}

#[get("/{project}/{version}/{build}/download")]
pub async fn download_build(path: Path<(String, String, String)>) -> Result<HttpResponse> {
    let (project, version, build) = path.into_inner();

    todo!()
}
