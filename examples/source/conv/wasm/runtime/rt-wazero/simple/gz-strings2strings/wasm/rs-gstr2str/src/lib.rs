#[cfg(feature = "smaller_wasm_wee_alloc")]
#[global_allocator]
static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;

use std::io::Read;

use flate2::read::GzDecoder;

static mut INPUT_GZIP_BYTES: Vec<u8> = vec![];

static mut OUTPUT_DECODED_BYTES: Vec<u8> = vec![];

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn resize_i(sz: i32) -> i32 {
    let mi: &mut Vec<u8> = unsafe { &mut INPUT_GZIP_BYTES };
    mi.resize(sz as usize, 0);
    mi.capacity().try_into().ok().unwrap_or(-1)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn reset_o(sz: i32) -> i32 {
    let mo: &mut Vec<u8> = unsafe { &mut OUTPUT_DECODED_BYTES };
    let cap: usize = mo.capacity();
    let add: usize = (sz as usize).saturating_sub(cap);
    mo.try_reserve(add)
        .ok()
        .and_then(|_| mo.capacity().try_into().ok())
        .unwrap_or(-1)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn offset_i() -> *mut u8 {
    let mi: &mut Vec<u8> = unsafe { &mut INPUT_GZIP_BYTES };
    mi.as_mut_ptr()
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn offset_o() -> *mut u8 {
    let mo: &mut Vec<u8> = unsafe { &mut OUTPUT_DECODED_BYTES };
    mo.as_mut_ptr()
}

fn rdr2buf<R>(rdr: R, buf: &mut Vec<u8>) -> Result<usize, &'static str>
where
    R: Read,
{
    buf.clear();
    let mut dec: GzDecoder<_> = GzDecoder::new(rdr);
    dec.read_to_end(buf).map_err(|_| "unable to read gz bytes")
}

fn slice2buf(s: &[u8], buf: &mut Vec<u8>) -> Result<usize, &'static str> {
    rdr2buf(s, buf)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn convert() -> i32 {
    let i: &Vec<u8> = unsafe { &INPUT_GZIP_BYTES };
    let o: &mut Vec<u8> = unsafe { &mut OUTPUT_DECODED_BYTES };
    slice2buf(i, o)
        .ok()
        .and_then(|u| u.try_into().ok())
        .unwrap_or(-1)
}
