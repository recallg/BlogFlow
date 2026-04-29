import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import "./style.css";
import "github-markdown-css/github-markdown.css";

const app = createApp(App);
app.use(router);
app.mount("#app");
