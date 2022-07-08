use actix_web::web::ServiceConfig;

mod builds;
mod projects;
mod versions;

// todo: fix duplicated code in router/ later
pub fn setup(config: &mut ServiceConfig) {
    config.configure(projects::routes);
    config.configure(versions::routes);
    config.configure(builds::routes);
}
