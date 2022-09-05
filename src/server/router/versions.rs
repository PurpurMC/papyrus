use crate::types::response::{VersionResponse, VersionResponseDetailed};
use crate::types::Response;
use actix_web::get;
use actix_web::web::{Path, ServiceConfig};

pub fn routes(config: &mut ServiceConfig) {
    config.service(get_version);
    config.service(get_version_detailed);
}

#[get("/{project}/{version}")]
pub async fn get_version(path: Path<(String, String)>) -> Response<VersionResponse> {
    let (project, version) = path.into_inner();

    todo!()
}

#[get("/{project}/{version}/detailed")]
pub async fn get_version_detailed(
    path: Path<(String, String)>,
) -> Response<VersionResponseDetailed> {
    let (project, version) = path.into_inner();

    todo!()
}
