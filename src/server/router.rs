use actix_web::web::ServiceConfig;

mod builds;
mod create;
mod projects;
mod versions;

pub fn setup(config: &mut ServiceConfig) {
    config.configure(projects::routes);
    config.configure(versions::routes);
    config.configure(builds::routes);
    config.configure(create::routes);
}
