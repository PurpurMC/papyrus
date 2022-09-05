use actix_multipart_extract::{File, MultipartForm};
use serde::Deserialize;

#[derive(Deserialize, MultipartForm)]
pub struct UploadFile {
    #[multipart(max_size = 100_000_000)] // 100MB
    pub file: File,
    pub extension: Option<String>,
}

#[derive(Deserialize)]
pub struct CreateBuild {
    pub file_id: String,
    pub project: String,
    pub version: String,
    pub build: String,
    pub result: String,
    pub duration: u64,
    pub timestamp: u64,
}
