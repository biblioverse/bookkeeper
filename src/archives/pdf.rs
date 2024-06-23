use super::BookInfo;
use archive_reader::error::Result;
use lopdf::{Document, Object};
use std::path::Path;

pub fn get_book_info_pdf(file: &Path) -> Result<BookInfo, Box<dyn std::error::Error>> {
    let doc = Document::load(file)?;
    let pages = doc.get_pages();

    let mut title: Option<String> = file.file_stem().map(|a| format!("{:?}", a));
    let mut authors: Option<Vec<String>> = None;
    let mut keywords: Option<Vec<String>> = None;

    if let Ok(info_obj_id) = doc.trailer.get(b"Info").and_then(Object::as_reference) {
        if let Ok(info_dict) = doc.get_object(info_obj_id).and_then(Object::as_dict) {
            if let Ok(raw_value) = info_dict.get(b"Author").and_then(Object::as_str) {
                let value = String::from_utf8_lossy(raw_value).to_string();

                if !value.is_empty() {
                    authors = Some(vec![value]);
                }
            }

            if let Ok(raw_value) = info_dict.get(b"Title").and_then(Object::as_str) {
                let value = String::from_utf8_lossy(raw_value).to_string();

                if !value.is_empty() && value != ".pdf" {
                    title = Some(value);
                }
            }

            if let Ok(raw_value) = info_dict.get(b"Keywords").and_then(Object::as_str) {
                let value = String::from_utf8_lossy(raw_value).to_string();

                if !value.is_empty() {
                    keywords = Some(
                        value
                            .split(',')
                            .map(|s| s.trim().to_string())
                            .filter(|s| !s.is_empty())
                            .collect(),
                    );
                }
            }

            // Print all metadata
            //for (key, value) in info_dict.iter() {
            //    println!("METADATA: [{}] [{:?}]", String::from_utf8(key.to_vec())?, value);
            //}
        }
    }

    Ok(BookInfo {
        title: title.unwrap_or("".to_string()),
        pages: pages.len(),
        authors,
        keywords,
        publisher: None,
        published_date: None,
    })
}

#[test]
fn read_metadata() -> Result<(), Box<dyn std::error::Error>> {
    let book = get_book_info_pdf(Path::new("fixtures/testfile.pdf"))?;
    assert_eq!(book.pages, 1);
    assert_eq!(book.title, "Title of the Book".to_string());
    assert_eq!(book.authors, Some(vec!["The Author".to_string()]));
    assert_eq!(
        book.keywords,
        Some(vec!["book".to_string(), "fantasy".to_string()])
    );

    Ok(())
}
