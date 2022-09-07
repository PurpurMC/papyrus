use crate::types::models::{Build, Project, Version};
use crate::types::response::{
    VersionResponse, VersionResponseBuilds, VersionResponseDetailed, VersionResponseDetailedBuilds,
};
use crate::types::Result;
use crate::utils::{response, verify};
use crate::SqlitePool;
use actix_web::web::{Data, Path, ServiceConfig};
use actix_web::{get, HttpResponse};

pub fn routes(config: &mut ServiceConfig) {
    config.service(get_version);
    config.service(get_version_detailed);
}

#[get("/{project}/{version}")]
pub async fn get_version(
    pool: Data<SqlitePool>,
    path: Path<(String, String)>,
) -> Result<HttpResponse> {
    let (project, version) = path.into_inner();
    let (latest, builds) = get_version_info(&project, &version, &pool).await?;

    let latest = latest.map(|build| build.name.clone());
    let builds = builds.iter().map(|build| build.name.clone()).collect();

    response(VersionResponse {
        project,
        version,
        builds: VersionResponseBuilds {
            latest,
            all: builds,
        },
    })
}

#[get("/{project}/{version}/detailed")]
pub async fn get_version_detailed(
    pool: Data<SqlitePool>,
    path: Path<(String, String)>,
) -> Result<HttpResponse> {
    let (project, version) = path.into_inner();
    let (latest, builds) = get_version_info(&project, &version, &pool).await?;

    let latest = match latest {
        Some(build) => Some(build.to_response(&project, &version, &pool).await?),
        None => None,
    };

    let mut build_responses = Vec::new();
    for build in builds {
        build_responses.push(build.to_response(&project, &version, &pool).await?)
    }

    response(VersionResponseDetailed {
        project,
        version,
        builds: VersionResponseDetailedBuilds {
            latest,
            all: build_responses,
        },
    })
}

async fn get_version_info(
    project: &String,
    version: &String,
    pool: &Data<SqlitePool>,
) -> Result<(Option<Build>, Vec<Build>)> {
    let project = verify(Project::find_one(&project, &pool).await?)?;
    let version = verify(Version::find_one(&project.id, &version, &pool).await?)?;

    let mut builds = Build::find_all(&version.id, &pool).await?;
    builds.sort_by(|a, b| b.created_at.cmp(&a.created_at));

    let latest: Option<Build> = {
        let mut found: Option<Build> = None;

        for build in builds.clone() {
            if build.result == "SUCCESS" {
                found = Some(build);
                break;
            }
        }

        found
    };

    Ok((latest, builds))
}
