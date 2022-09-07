use crate::types::models::{Build, File, Project, Version};
use crate::types::{Error, Result};
use crate::utils::{response, verify};
use crate::{Config, SqlitePool};
use actix_files::NamedFile;
use actix_web::http::header::{ContentDisposition, DispositionParam, DispositionType};
use actix_web::web::{Data, Path, ServiceConfig};
use actix_web::{get, HttpRequest, HttpResponse};
use mime_guess::Mime;
use std::path;
use std::str::FromStr;

pub fn routes(config: &mut ServiceConfig) {
    config.service(get_build);
    config.service(download_build);
}

#[get("/{project}/{version}/{build}")]
pub async fn get_build(
    pool: Data<SqlitePool>,
    path: Path<(String, String, String)>,
) -> Result<HttpResponse> {
    let (project, version, build) = path.into_inner();
    let project = verify(Project::find_one(&project, &pool).await?)?;
    let version = verify(Version::find_one(&project.id, &version, &pool).await?)?;
    let build = verify(Build::find_one(&version.id, &build, &pool).await?)?;

    response(
        build
            .to_response(&project.name, &version.name, &pool)
            .await?,
    )
}

#[get("/{project}/{version}/{build}/download")]
pub async fn download_build(
    request: HttpRequest,
    pool: Data<SqlitePool>,
    config: Data<Config>,
    path: Path<(String, String, String)>,
) -> Result<HttpResponse> {
    let (project, version, build) = path.into_inner();
    let project = verify(Project::find_one(&project, &pool).await?)?;
    let version = verify(Version::find_one(&project.id, &version, &pool).await?)?;
    let build = verify(Build::find_one(&version.id, &build, &pool).await?)?;
    let file = verify(File::find_one(&build.id, &pool).await?)?;

    let path = format!("{}/files/{}", &config.database, &file.id);
    let path = path::Path::new(&path);
    if !path.exists() {
        return Err(Error::NotFound);
    }

    let file_extension = file.extension;
    let file = NamedFile::open_async(path).await?;

    Ok(file
        .set_content_type(
            mime_guess::from_ext(&file_extension)
                .first()
                .unwrap_or(Mime::from_str("application/octet-stream").unwrap()),
        )
        .set_content_disposition(ContentDisposition {
            disposition: DispositionType::Attachment,
            parameters: vec![DispositionParam::Filename(format!(
                "{}-{}-{}.{}",
                &project.name, &version.name, &build.name, &file_extension
            ))],
        })
        .into_response(&request))
}
