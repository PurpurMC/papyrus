use crate::types::Result;
use crate::SqlitePool;

mod sqlite;

pub use sqlite::*;

trait Model<Params> {
    fn all(pool: &SqlitePool) -> Result<Vec<Box<Self>>>;

    fn find_by(params: Params, pool: &SqlitePool) -> Result<Option<Box<Self>>>;
}
