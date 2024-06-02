use std::ffi::OsStr;
use std::path::Path;

pub fn is_zip(entry: &Path) -> bool {
    let cbr: &OsStr = OsStr::new("cbr");
    let cbz: &OsStr = OsStr::new("cbz");

    if let Some(ext) = entry.extension() {
        let lowercase_ext = ext.to_ascii_lowercase();

        if lowercase_ext == cbr || lowercase_ext == cbz {
            return true;
        }
    }

    return false;
}

pub fn is_pdf(entry: &Path) -> bool {
    let pdf: &OsStr = OsStr::new("pdf");

    if let Some(ext) = entry.extension() {
        let lowercase_ext = ext.to_ascii_lowercase();

        if lowercase_ext == pdf {
            return true;
        }
    }

    return false;
}
