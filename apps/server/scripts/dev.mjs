import { spawn, spawnSync } from "node:child_process";
import { join } from "node:path";

const rootDir = join(import.meta.dirname, "..");

const startDependencies = () => {
  const result = spawnSync(
    "docker",
    ["compose", "up", "postgres", "redis", "-d", "--wait"],
    {
      cwd: rootDir,
      stdio: "inherit",
    }
  );

  if (result.status !== 0) {
    console.error("Failed to start postgres/redis. Is Docker Desktop running?");
    process.exit(result.status ?? 1);
  }
};

const main = () => {
  console.log("Starting postgres and redis...");
  startDependencies();

  const server = spawn(
    process.execPath,
    [join(rootDir, "scripts/with-go-env.mjs"), "run", "./cmd/server"],
    {
      cwd: rootDir,
      stdio: "inherit",
    }
  );

  const handleExit = function handleExit(code) {
    process.exit(code ?? 1);
  };

  server.on("exit", handleExit);
};

main();
