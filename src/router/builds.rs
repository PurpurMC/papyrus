use crate::models::{Build, Commit};
use crate::{utils, Config, SqlitePool};
use actix_files::NamedFile;
use actix_web::http::header::{ContentDisposition, DispositionParam, DispositionType};
use actix_web::web::{Data, Path, ServiceConfig};
use actix_web::{get, HttpRequest, HttpResponse};
use mime_guess::Mime;
use serde::Serialize;
use serde_json::json;
use std::str::FromStr;

pub fn routes(config: &mut ServiceConfig) {
    config.service(get_build);
    config.service(download_build);
}

#[get("/{project}/{version}/{build}")]
async fn get_build(pool: Data<SqlitePool>, path: Path<(String, String, String)>) -> HttpResponse {
    let (project, version, build) = path.into_inner();

    let project = match utils::router::project(&pool, &project).await {
        Ok(projects) => projects,
        Err(err) => return err,
    };

    let version = match utils::router::version(&pool, &project, &version).await {
        Ok(version) => version,
        Err(err) => return err,
    };

    let build = {
        if build == "latest" {
            let builds = match Build::all(&pool, &version.id).await {
                Ok(builds) => builds,
                Err(err) => {
                    return HttpResponse::InternalServerError()
                        .json(json!({ "error": err.to_string() }))
                }
            };

            match builds
                .iter()
                .filter(|build| build.hash.is_some())
                .rev()
                .find(|build| build.result != "FAILURE")
            {
                Some(build) => build.clone(),
                None => {
                    return HttpResponse::NotFound().json(json!({ "error": "Build not found" }))
                }
            }
        } else {
            match utils::router::build(&pool, &version, &build).await {
                Ok(build) => build,
                Err(err) => return err,
            }
        }
    };

    if build.hash.is_none() {
        return HttpResponse::NotFound().json(json!({ "error": "Build not found" }));
    }

    HttpResponse::Ok().json(BuildResponse {
        project: project.name.clone(),
        version: version.name.clone(),
        build: build.name.clone(),
        result: build.result.clone(),
        commits: match commits(&pool, &build).await {
            Ok(commits) => commits,
            Err(err) => return err,
        },
        md5: build.hash.unwrap_or("".into()),
        duration: build.duration,
        timestamp: build.timestamp,
    })
}

#[get("/{project}/{version}/{build}/download")]
async fn download_build(
    req: HttpRequest,
    pool: Data<SqlitePool>,
    config: Data<Config>,
    path: Path<(String, String, String)>,
) -> HttpResponse {
    let (project, version, build) = path.into_inner();
    let project = match utils::router::project(&pool, &project).await {
        Ok(projects) => projects,
        Err(err) => return err,
    };

    let version = match utils::router::version(&pool, &project, &version).await {
        Ok(version) => version,
        Err(err) => return err,
    };

    let build = match utils::router::build(&pool, &version, &build).await {
        Ok(build) => build,
        Err(err) => return err,
    };

    let file = match NamedFile::open_async(format!("{}/{}", config.database, build.id)).await {
        Ok(file) => file,
        Err(err) => {
            return HttpResponse::InternalServerError().json(json!({"error": err.to_string()}))
        }
    };

    let file_extension = build.file_extension.unwrap_or("".into());
    file.set_content_type(
        mime_guess::from_ext(&file_extension)
            .first()
            .unwrap_or(Mime::from_str("application/octet-stream").unwrap()),
    )
    .set_content_disposition(ContentDisposition {
        disposition: DispositionType::Attachment,
        parameters: vec![DispositionParam::Filename(format!(
            "{}-{}-{}.{}",
            project.name, version.name, build.name, file_extension
        ))],
    })
    .into_response(&req)
}

#[derive(Serialize)]
pub struct BuildResponse {
    pub project: String,
    pub version: String,
    pub build: String,
    pub result: String,
    pub commits: Vec<CommitResponse>,
    pub md5: String,
    pub duration: i64,
    pub timestamp: i64,
}

#[derive(Serialize)]
pub struct CommitResponse {
    pub author: String,
    pub email: String,
    pub description: String,
    pub hash: String,
    pub timestamp: i64,
}

pub async fn commits(
    pool: &SqlitePool,
    build: &Build,
) -> Result<Vec<CommitResponse>, HttpResponse> {
    let commits = match Commit::get(pool, &build.id).await {
        Ok(commits) => commits,
        Err(err) => {
            return Err(
                HttpResponse::InternalServerError().json(json!({ "error": err.to_string() }))
            )
        }
    };

    Ok(commits
        .iter()
        .map(|commit| CommitResponse {
            author: commit.author.clone(),
            email: commit.email.clone(),
            description: commit.description.clone(),
            hash: commit.hash.clone(),
            timestamp: commit.timestamp,
        })
        .collect::<Vec<CommitResponse>>())
}
