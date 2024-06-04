use super::BookInfo;
use archive_reader::error::Result;
use lopdf::Document;
use lopdf::Object::String as PdfString;
use std::ffi::OsStr;
use std::path::Path;

pub fn get_book_info_pdf(file: &Path) -> Result<BookInfo, Box<dyn std::error::Error>> {
    let doc = Document::load(file)?;
    let pages = doc.get_pages();

    // TODO :: set default to null
    let mut authors: Option<Vec<String>> = None;
    if let Ok(info) = doc.trailer.get(b"Info") {
        if let Ok(info_dict) = info.as_dict() {
            if let Ok(PdfString(raw_author, _)) = info_dict.get(b"Author") {
                authors = Some(vec![format!("{:?}", String::from_utf8_lossy(raw_author))]);
            }

            // TODO : get other metadata
        }
    }

    Ok(BookInfo {
        title: format!("{:?}", file.file_stem().unwrap_or(OsStr::new("stuff"))),
        pages: pages.len(),
        authors,
        publisher: None,
        published_date: None,
        keywords: None,
    })
}
