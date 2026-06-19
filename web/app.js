"use strict";

const playButton = document.querySelector("#play");
const statusText = document.querySelector("#status");
const launcher = document.querySelector("#launcher");

if (!WebAssembly.instantiateStreaming) {
  WebAssembly.instantiateStreaming = async (response, importObject) => {
    const source = await (await response).arrayBuffer();
    return WebAssembly.instantiate(source, importObject);
  };
}

const go = new Go();
let instance;

WebAssembly.instantiateStreaming(fetch("gong.wasm"), go.importObject)
  .then((result) => {
    instance = result.instance;
    statusText.textContent = "Ready.";
    playButton.disabled = false;
  })
  .catch((error) => {
    console.error(error);
    statusText.textContent = "The game could not be loaded. Please reload the page.";
  });

playButton.addEventListener("click", () => {
  playButton.disabled = true;
  statusText.textContent = "Starting…";
  launcher.hidden = true;

  go.run(instance).catch((error) => {
    console.error(error);
    launcher.hidden = false;
    statusText.textContent = "The game stopped unexpectedly. Reload to try again.";
  });
});
