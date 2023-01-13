# Eiger Labs Programming Test

This is the output of a programming test set by Eiger Labs for the position of Senior Software Engineer.

- Files are deterministically chunked via Rabin fingerprinting (rolling digest).
- Fast xxHash fingerprints are generated for each chunk for both the new file and old file.
- Fingerprints are compared, and if they are different between old file chunk and new file chunk, the new file chunk byte slice, position, length and chunk number is stored.

The following test cases were used:

- alpha: Simple alphabet with replaced character
- duplication: Large file repeated twice
- prepend: Large file with string inserted into the beginning of the file
- insertion: Small lorem ipsum with text inserted in the middle
- removal: Small lorem ipsum with a sentence removed

The patch functionality is not yet implemented.

## References:

- Splitting Data with Content-Defined Chunking: https://blog.gopheracademy.com/advent-2018/split-data-with-cdc/
- Foundation - Introducing Content Defined Chunking (CDC): https://restic.net/blog/2015-09-12/restic-foundation1-cdc/
- Rabin fingerprint: https://en.wikipedia.org/wiki/Rabin_fingerprint

### Rdiff Help Text

You can use rdiff to update files, much like rsync does. However, unlike rsync, rdiff puts you in control. There are three steps to updating a file: signature, delta, and patch.

- Use the signature subcommand to generate a small signature-file from the old-file.
- Use the delta subcommand to generate a small delta-file from the signature-file to the new-file.
- Use the patch subcommand to apply the delta-file to the old-file to regenerate the new-file.

### Spec v4 (2021-03-09)

Make a rolling hash based file diffing algorithm. When comparing original and an updated version of an input,
it should return a description ("delta") which can be used to upgrade an original version of the file into the
new file. The description contains the chunks which:

- Can be reused from the original file
- have been added or modified and thus would need to be synchronized

The real-world use case for this type of construct could be a distributed file storage system.

This reduces the need for bandwidth and storage. If many people have the same file stored on Dropbox, for example, there's no need to upload it again.

A library that does a similar thing is rdiff. You don't need to fulfill the patch part of the API, only signature and delta.

**Requirements:**

- Hashing function gets the data as a parameter. Separate possible filesystem operations.
- Chunk size can be fixed or dynamic, but must be split to at least two chunks on any sufficiently sized data.
- Should be able to recognize changes between chunks. Only the exact differing locations should be added to the delta.
- Well-written unit tests function well in describing the operation, no UI necessary.

**Checklist:**

- Input/output operations are separated from the calculations
- detects chunk changes and/or additions
- detects chunk removals
- detects additions between chunks with shifted original chunks
