use serde::Serialize;

#[derive(Serialize)]
pub struct DefaultResponse {
    pub message: String,
}

#[derive(Serialize)]
pub struct ProjectsResponse {
    pub projects: Vec<String>,
}

#[derive(Serialize)]
pub struct ProjectResponse {
    pub project: String,
    pub versions: Vec<String>,
}

#[derive(Serialize)]
pub struct VersionResponse {
    pub project: String,
    pub version: String,
    pub builds: VersionResponseBuilds,
}

#[derive(Serialize)]
pub struct VersionResponseBuilds {
    pub latest: Option<String>,
    pub all: Vec<String>,
}

#[derive(Serialize)]
pub struct VersionResponseDetailed {
    pub project: String,
    pub version: String,
    pub builds: VersionResponseBuilds,
}

#[derive(Serialize)]
pub struct VersionResponseDetailedBuilds {
    pub latest: Option<BuildResponse>,
    pub all: Vec<BuildResponse>,
}

#[derive(Serialize)]
pub struct BuildResponse {
    pub project: String,
    pub version: String,
    pub build: String,
    pub commits: Vec<BuildResponseCommit>,
    pub result: String,
    pub md5: String,
    pub duration: u64,
    pub timestamp: u64,
}

#[derive(Serialize)]
pub struct BuildResponseCommit {
    pub author: String,
    pub email: String,
    pub description: String,
    pub hash: String,
    pub timestamp: u64,
}

#[derive(Serialize)]
pub struct UploadFileResponse {
    pub id: String,
}
