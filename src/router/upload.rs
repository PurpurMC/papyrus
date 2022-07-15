use crate::middleware::Authentication;
use crate::models::{Project, Version};
use crate::Config;
use actix_multipart_extract::{Multipart, MultipartForm};
use actix_web::web::{Data, Json, ServiceConfig};
use actix_web::{post, web, HttpResponse};
use nanoid::nanoid;
use serde::Deserialize;
use serde_json::json;
use sqlx::SqlitePool;
use std::fs::File;
use std::io::Write;
use std::path::Path;

// todo: clean up after ourselves if something goes wrong
pub fn routes(config: &mut ServiceConfig) {
    config.service(
        web::scope("/upload")
            .wrap(Authentication)
            .service(create_build)
            .service(upload_file),
    );
}

#[derive(Deserialize)]
struct CreatePayload {
    project: String,
    version: String,
    build: String,
    result: String,
    commits: Vec<CreatePayloadCommit>,
    duration: i64,
    timestamp: i64,
}

#[derive(Deserialize)]
struct CreatePayloadCommit {
    author: String,
    email: String,
    description: String,
    hash: String,
    timestamp: i64,
}

#[post("/create")]
async fn create_build(pool: Data<SqlitePool>, payload: Json<CreatePayload>) -> HttpResponse {
    let version_id = match get_version_id(&pool, &payload.project, &payload.version).await {
        Ok(id) => id,
        Err(err) => return err,
    };

    let build = payload.build.clone();
    let optional = match sqlx::query!(
        "SELECT id FROM builds WHERE name = ? AND version_id = ?",
        build,
        version_id
    )
    .fetch_optional(pool.as_ref())
    .await
    {
        Ok(optional) => optional,
        Err(err) => {
            return HttpResponse::InternalServerError().json(json!({ "error": err.to_string() }))
        }
    };

    if optional.is_some() {
        return HttpResponse::Conflict().json(json!({ "error": "Build already exists" }));
    }

    let build_id = nanoid!();
    match sqlx::query!("INSERT INTO builds (id, name, version_id, result, duration, timestamp) VALUES (?, ?, ?, ?, ?, ?)", build_id, build, version_id, payload.result, payload.duration, payload.timestamp).execute(pool.as_ref()).await {
        Ok(_) => (),
        Err(err) => return HttpResponse::InternalServerError().json(json!({ "error": err.to_string() })),
    };

    for commit in &payload.commits {
        let id = nanoid!();
        match sqlx::query!("INSERT INTO commits (id, build_id, author, email, description, hash, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?)", id, build_id, commit.author, commit.email, commit.description, commit.hash, commit.timestamp).execute(pool.as_ref()).await {
            Ok(_) => (),
            Err(err) => return HttpResponse::InternalServerError().json(json!({ "error": err.to_string() })),
        };
    }

    HttpResponse::Ok().json(json!({
        "status": "ok",
        "build_id": build_id,
    }))
}

async fn get_version_id(
    pool: &SqlitePool,
    project: &String,
    version: &String,
) -> Result<String, HttpResponse> {
    let id = nanoid!();
    match sqlx::query!(
        "INSERT OR IGNORE INTO projects (id, name) VALUES (?, ?)",
        id,
        project
    )
    .execute(pool)
    .await
    {
        Ok(_) => (),
        Err(err) => {
            return Err(
                HttpResponse::InternalServerError().json(json!({ "error": err.to_string() }))
            )
        }
    };

    let project = match Project::get(&pool, &project).await {
        Ok(project) => project,
        Err(err) => {
            return Err(
                HttpResponse::InternalServerError().json(json!({ "error": err.to_string() }))
            )
        }
    }
    .unwrap();

    let id = nanoid!();
    match sqlx::query!(
        "INSERT OR IGNORE INTO versions (id, name, project_id) VALUES (?, ?, ?)",
        id,
        version,
        project.id
    )
    .execute(pool)
    .await
    {
        Ok(_) => (),
        Err(err) => {
            return Err(
                HttpResponse::InternalServerError().json(json!({ "error": err.to_string() }))
            )
        }
    };

    let version = match Version::get(&pool, &version, &project.id).await {
        Ok(version) => version,
        Err(err) => {
            return Err(
                HttpResponse::InternalServerError().json(json!({ "error": err.to_string() }))
            )
        }
    }
    .unwrap();

    Ok(version.id)
}

#[derive(Deserialize, MultipartForm)]
struct UploadPayload {
    build_id: String,
    #[multipart(max_size = 100_000_000)] // 100MB
    file: actix_multipart_extract::File,
    file_extension: String,
}

#[post("/file")]
async fn upload_file(
    pool: Data<SqlitePool>,
    config: Data<Config>,
    payload: Multipart<UploadPayload>,
) -> HttpResponse {
    let build_id = payload.build_id.clone();
    let optional = match sqlx::query!("SELECT id FROM builds WHERE id = ?", build_id)
        .fetch_optional(pool.as_ref())
        .await
    {
        Ok(optional) => optional,
        Err(err) => {
            return HttpResponse::InternalServerError().json(json!({ "error": err.to_string() }))
        }
    };

    if optional.is_none() {
        return HttpResponse::NotFound().json(json!({ "error": "Build not found" }));
    }

    let filename = format!("{}/{}", &config.database, build_id);
    let path = Path::new(&filename);
    if path.exists() {
        return HttpResponse::Conflict().json(json!({ "error": "Build already exists" }));
    }

    let mut file = match File::create(&filename) {
        Ok(file) => file,
        Err(err) => {
            return HttpResponse::InternalServerError().json(json!({ "error": err.to_string() }))
        }
    };

    match file.write_all(&payload.file.bytes) {
        Ok(_) => (),
        Err(err) => {
            return HttpResponse::InternalServerError().json(json!({ "error": err.to_string() }))
        }
    };

    let md5 = format!("{:x}", md5::compute(&payload.file.bytes));
    match sqlx::query!(
        "UPDATE builds SET hash = ?, file_extension = ? WHERE id = ?",
        md5,
        payload.file_extension,
        build_id
    )
    .execute(pool.as_ref())
    .await
    {
        Ok(_) => (),
        Err(err) => {
            return HttpResponse::InternalServerError().json(json!({ "error": err.to_string() }))
        }
    };

    HttpResponse::Ok().json(json!({
        "status": "ok",
    }))
}
