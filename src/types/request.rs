use actix_multipart_extract::{File, MultipartForm};
use serde::Deserialize;

#[derive(Deserialize, MultipartForm)]
pub struct UploadFileRequest {
    #[multipart(max_size = 100_000_000)] // 100MB
    pub file: File,
    pub extension: Option<String>,
}

#[derive(Deserialize)]
pub struct CreateBuildRequest {
    pub file_id: String,
    pub project: String,
    pub version: String,
    pub build: String,
    pub commits: Vec<CreateBuildRequestCommit>,
    pub result: String,
    pub duration: i64,
    pub timestamp: i64,
}

#[derive(Deserialize)]
pub struct CreateBuildRequestCommit {
    pub author: String,
    pub email: String,
    pub description: String,
    pub hash: String,
    pub timestamp: i64,
}
