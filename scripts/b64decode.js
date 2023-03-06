#!/usr/bin/env node

const fs = require('fs');
const data = fs.readFileSync(0, 'utf-8');
function fromBinary(encoded) {
  const binary = atob(encoded);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < bytes.length; i++) {
    bytes[i] = binary.charCodeAt(i);
  }
  return String.fromCharCode(...new Uint16Array(bytes.buffer));
}
process.stdout.write(fromBinary(data)+"\n");
