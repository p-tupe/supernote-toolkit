## v0.4 - Wed Feb 11 16:39:07 EST 2026

- Add support for Nomad X2 devices.

- Add support for Horizontal orientations.

- Add cli flag (-device) and UI radio buttons (Select Device) to manually impose a device for all notes for the conversion.

## v0.3 - Tue Feb 10 16:31:47 EST 2026

- Skip conversion when .note file is older than the existing .png/.pdf file; add a new "force" option to convert all files regardless.

- Recurse into sub-folders by default.

- Converted PNGs now go into folder without the .note extention folders since it conflicted when output was same as input.
