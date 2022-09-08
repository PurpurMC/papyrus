use crate::types::models::{Build, File, Project, Version};
use crate::types::response::{
    BuildResponse, VersionResponse, VersionResponseBuilds, VersionResponseDetailed,
    VersionResponseDetailedBuilds,
};
use crate::types::Result;
use crate::utils::{response, verify};
use crate::SqlitePool;
use actix_web::web::{Data, Path, ServiceConfig};
use actix_web::{get, HttpRequest, HttpResponse};
use qstring::QString;

pub fn routes(config: &mut ServiceConfig) {
    config.service(get_version);
}

#[get("/{project}/{version}")]
pub async fn get_version(
    request: HttpRequest,
    pool: Data<SqlitePool>,
    path: Path<(String, String)>,
) -> Result<HttpResponse> {
    let (project, version) = path.into_inner();
    let project = verify(Project::find_one(&project, &pool).await?)?;
    let version = verify(Version::find_one(&project.id, &version, &pool).await?)?;

    let mut all_builds = Build::find_all(&version.id, &pool).await?;
    all_builds.sort_by(|a, b| a.created_at.cmp(&b.created_at));

    let mut builds = Vec::new();
    for build in all_builds {
        if build.result == "FAILURE" || File::find_one(&build.id, &pool).await?.is_some() {
            builds.push(build);
        }
    }

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

    if QString::from(request.query_string()).has("detailed") {
        let latest = match latest {
            Some(build) => Some(
                build
                    .to_response(&project.name, &version.name, &pool)
                    .await?,
            ),
            None => None,
        };

        let latest: Option<BuildResponse> = match latest {
            Some(option) => option,
            None => None,
        };

        let mut build_responses = Vec::new();
        for build in builds {
            let response = build
                .to_response(&project.name, &version.name, &pool)
                .await?;
            if let Some(response) = response {
                build_responses.push(response);
            }
        }

        response(VersionResponseDetailed {
            project: project.name,
            version: version.name,
            builds: VersionResponseDetailedBuilds {
                latest,
                all: build_responses,
            },
        })
    } else {
        let latest = latest.map(|build| build.name.clone());
        let builds = builds.iter().map(|build| build.name.clone()).collect();

        response(VersionResponse {
            project: project.name,
            version: version.name,
            builds: VersionResponseBuilds {
                latest,
                all: builds,
            },
        })
    }
}
