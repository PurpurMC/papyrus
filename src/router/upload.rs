use std::fmt::format;
use std::fs::File;
use std::io::{Read, Write};
use std::path::Path;
use actix_multipart_extract::{Multipart, MultipartForm};
use actix_web::{HttpResponse, post};
use actix_web::web::{Data, Json, ServiceConfig};
use nanoid::nanoid;
use serde::Deserialize;
use serde_json::json;
use sqlx::SqlitePool;
use crate::Config;
use crate::models::{Build, Project, Version};

pub fn routes(config: &mut ServiceConfig) {
    config.service(create_build);
    config.service(upload_file);
}

#[derive(Deserialize)]
struct CreatePayload {
    project: String,
    version: String,
    build: String,
    result: String,
    commits: Vec<PayloadCommit>,
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

// todo: clean up after ourselves if something goes wrong
#[post("/upload/create")]
async fn create_build(pool: Data<SqlitePool>, payload: Json<CreatePayload>) -> HttpResponse {
    let version_id = match get_version_id(&pool, &payload.project, &payload.version).await {
        Ok(id) => id,
        Err(err) => return err,
    };

    let build = payload.build.clone();
    let version_id = version_id.clone();
    let optional = match sqlx::query!("SELECT id FROM builds WHERE name = ? AND version_id = ?", build, version_id).fetch_optional(pool.as_ref()).await {
        Ok(optional) => optional,
        Err(err) => return HttpResponse::InternalServerError().json(json!({ "error": err.to_string() })),
    };

    if optional.is_some() {
        return HttpResponse::Conflict().json(json!({ "error": "Build already exists" }));
    }

    let build_id = nanoid!(10);
    match sqlx::query!("INSERT INTO builds (id, name, version_id, result, duration, timestamp) VALUES (?, ?, ?, ?, ?, ?)", build_id, build, version_id, payload.result, payload.duration, payload.timestamp).execute(pool.as_ref()).await {
        Ok(_) => (),
        Err(err) => return HttpResponse::InternalServerError().json(json!({ "error": err.to_string() })),
    };

    for commit in &payload.commits {
        let id = nanoid!(10);
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

#[derive(Deserialize, MultipartForm)]
struct UploadPayload {
    build_id: String,
    #[multipart(max_size = 100_000_000)] // 100MB
    file: actix_multipart_extract::File,
    file_extension: String,
}

// todo: clean up after ourselves if something goes wrong
#[post("/upload/file")]
async fn upload_file(pool: Data<SqlitePool>, config: Data<Config>, payload: Multipart<UploadPayload>) -> HttpResponse {
    let build_id = payload.build_id.clone();
    let optional = match sqlx::query!("SELECT id FROM builds WHERE id = ?", build_id).fetch_optional(pool.as_ref()).await {
        Ok(optional) => optional,
        Err(err) => return HttpResponse::InternalServerError().json(json!({ "error": err.to_string() })),
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
        Err(err) => return HttpResponse::InternalServerError().json(json!({ "error": err.to_string() })),
    };

    match file.write_all(&payload.file.bytes) {
        Ok(_) => (),
        Err(err) => return HttpResponse::InternalServerError().json(json!({ "error": err.to_string() })),
    };

    let md5 = format!("{:x}", md5::compute(&payload.file.bytes));
    sqlx::query!("UPDATE builds SET hash = ?, file_extension = ?, uploaded = TRUE WHERE id = ?", md5, payload.file_extension, build_id).execute(pool.as_ref()).await.unwrap();

    HttpResponse::Ok().json(json!({
        "status": "ok",
    }))
}

async fn get_version_id(pool: &SqlitePool, project: &String, version: &String) -> Result<String, HttpResponse> {
    let id = nanoid!(10);
    sqlx::query!("INSERT OR IGNORE INTO projects (id, name) VALUES (?, ?)", id, project).execute(pool).await;

    let project = match Project::get(&pool, &project).await {
        Ok(project) => project,
        Err(err) => return Err(HttpResponse::InternalServerError().json(json!({ "error": err.to_string() }))),
    }.unwrap();

    let id = nanoid!(10);
    sqlx::query!("INSERT OR IGNORE INTO versions (id, name, project_id) VALUES (?, ?, ?)", id, version, project.id).execute(pool).await;

    let version = match Version::get(&pool, &version, &project.id).await {
        Ok(version) => version,
        Err(err) => return Err(HttpResponse::InternalServerError().json(json!({ "error": err.to_string() }))),
    }.unwrap();

    Ok(version.id)
}
