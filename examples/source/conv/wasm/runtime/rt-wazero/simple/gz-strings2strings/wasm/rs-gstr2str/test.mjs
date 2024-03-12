import { readFile } from "node:fs/promises"

(() => {
    return Promise.resolve("rs_gstr2str.wasm")
    .then(name => readFile(name))
    .then(bytes => WebAssembly.instantiate(bytes))
    .then(pair => {
        const {
            module,
            instance,
        } = pair || {}
        const {
            memory,

            resize_i,
            reset_o,

            offset_i,
            offset_o,

            convert,
        } = instance?.exports || {}

        const gzipBytes = Buffer.from(
            "H4sIAAAAAAAEAzPUMdIx1jEBAGeLu+wHAAAA",
            "base64",
        )

        const icap = resize_i(gzipBytes.length)
        const ocap = reset_o(gzipBytes.length)

        const iptr = offset_i()
        const optr = offset_o()

        const iview = new Uint8Array(memory?.buffer, iptr, gzipBytes.length)
        iview.set(gzipBytes)

        const osz = convert()
        const oview = new Uint8Array(memory?.buffer, optr, osz)
        const dec = new TextDecoder()

        return dec.decode(oview)
    })
    .then(console.info)
    .catch(console.warn)
})()
