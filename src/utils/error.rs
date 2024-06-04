use std::error::Error;
use std::fmt;

#[derive(Debug)]
pub struct IllegalArgumentError {
    message: String,
}

impl IllegalArgumentError {
    pub fn new(message: String) -> Self {
        IllegalArgumentError { message }
    }
}

impl Error for IllegalArgumentError {}

impl fmt::Display for IllegalArgumentError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "Illegal argument: {}", self.message)
    }
}
