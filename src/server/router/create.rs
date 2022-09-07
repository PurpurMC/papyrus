use crate::server::middleware::RequireKey;
use crate::types::models::{Build, Commit, File, Project, TempFile, Version};
use crate::types::request::{CreateBuildRequest, UploadFileRequest};
use crate::types::response::{DefaultResponse, UploadFileResponse};
use crate::types::{Error, Result};
use crate::utils::{response, verify};
use crate::{Config, SqlitePool};
use actix_multipart_extract::Multipart;
use actix_web::web;
use actix_web::web::{Data, Json, ServiceConfig};
use actix_web::{post, HttpResponse};
use std::fs;
use std::io::Write;
use std::path::Path;

pub fn routes(config: &mut ServiceConfig) {
    config.service(
        web::scope("")
            .wrap(RequireKey)
            .service(upload_file)
            .service(create_build),
    );
}

#[post("/create/file")]
pub async fn upload_file(
    payload: Multipart<UploadFileRequest>,
    pool: Data<SqlitePool>,
    config: Data<Config>,
) -> Result<HttpResponse> {
    let mut temp_file = TempFile::create(&payload.extension, &pool).await?;

    loop {
        let path = format!("{}/temp/{}", &config.database, &temp_file.id);
        let path = Path::new(&path);
        if !path.exists() {
            let mut file = fs::File::create(path)?;
            file.write_all(&payload.file.bytes)?;

            break;
        }

        TempFile::delete(&temp_file.id, &pool).await?;
        temp_file = TempFile::create(&payload.extension, &pool).await?;
    }

    response(UploadFileResponse { id: temp_file.id })
}

#[post("/create/build")]
pub async fn create_build(
    payload: Json<CreateBuildRequest>,
    pool: Data<SqlitePool>,
    config: Data<Config>,
) -> Result<HttpResponse> {
    let temp_file = verify(TempFile::find_one(&payload.file_id, &pool).await?)?;

    let temp_path = format!("{}/temp/{}", &config.database, &temp_file.id);
    let temp_path = Path::new(&temp_path);
    if !temp_path.exists() {
        TempFile::delete(&temp_file.id, &pool).await?;
        return Err(Error::NotFound);
    }

    let project = Project::create(&payload.project, &pool).await?;
    let version = Version::create(&project.id, &payload.version, &pool).await?;
    let build: Build = match Build::find_one(&version.id, &payload.build, &pool).await? {
        Some(_) => return Err(Error::AlreadyExists),
        None => Ok::<Build, Error>(Build::create(&version.id, &payload, &pool).await?),
    }?;

    for commit in &payload.commits {
        Commit::create(&build.id, commit, &pool).await?;
    }

    let hash = md5::compute(fs::read(&temp_path)?);
    let hash = format!("{:x}", hash);

    let file = File::create(&build.id, &temp_file.extension, &hash, &pool).await?;
    let new_path = format!("{}/files/{}", &config.database, &file.id);
    match fs::rename(temp_path, new_path) {
        Ok(_) => Ok(()),
        Err(err) => {
            File::delete(&file.id, &pool).await?;
            Err(err)
        }
    }?;

    TempFile::delete(&temp_file.id, &pool).await?;

    response(DefaultResponse {
        message: "successfully created build".into(),
    })
}
