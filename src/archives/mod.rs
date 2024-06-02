use std::path::Path;
use std::error::Error;

pub mod cb;
pub mod pdf;

use cb::get_book_info_cba;
use pdf::get_book_info_pdf;
use super::utils::fileutils::{is_zip, is_pdf};
use super::utils::error::IllegalArgumentError;

pub struct BookInfo {
    pages: usize
}

impl BookInfo {
    pub fn page_count(&self) -> usize {
        return self.pages
    }
}

pub fn get_book_info(file: &Path) -> Result<BookInfo, Box<dyn Error>> {
    if is_zip(&file) {
        return Ok(get_book_info_cba(file)?);
    }

    if is_pdf(&file) {
        return Ok(get_book_info_pdf(file)?)
    }

    return Err(Box::new(IllegalArgumentError::new(format!("We don't know how to open this archive '{:?}'", file))));
}
