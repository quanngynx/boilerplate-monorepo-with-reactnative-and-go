import { spawnSync } from "node:child_process";
import { mkdirSync } from "node:fs";
import { join } from "node:path";

const gocache = join(process.cwd(), ".cache", "go");
mkdirSync(gocache, { recursive: true });

const [, , ...args] = process.argv;
const result = spawnSync("go", args, {
  env: { ...process.env, GOCACHE: gocache },
  stdio: "inherit",
});

process.exit(result.status ?? 1);
