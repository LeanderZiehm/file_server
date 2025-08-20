import crypto from "crypto";

export function generateFileId(filename: string): string {
  const timestamp = Date.now();
  return crypto.createHash("sha256").update(filename + timestamp).digest("hex");
}
