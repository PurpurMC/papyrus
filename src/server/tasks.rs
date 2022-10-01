use crate::types::Result;
use async_trait::async_trait;
use std::time::Duration;

pub use clear_temp::*;

mod clear_temp;

#[async_trait]
pub trait Task {
    fn duration(&self) -> Duration;
    async fn run(&self) -> Result<()>;
}
