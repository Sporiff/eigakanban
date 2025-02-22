<script setup lang="ts">
import { useRouter } from "vue-router";
import { useApiClientStore } from "@/stores/apiClient.ts";

const apiClientStore = useApiClientStore();
const router = useRouter();

let credential = "";
let password = "";
let success = true;
let status = "";

async function handleLogin() {
  try {
    const result = await apiClientStore.login(credential, password);

    success = result.success;
    status = result.status;

    if (result.success) {
      // Redirect after successful login
      await router.push("/"); // Navigate to home page
    }
  } catch (e: any) {
    success = false;
    status = e.message;
    credential = "";
    password = "";
  }
}
</script>

<template>
  <form @submit.prevent="handleLogin">
    <div class="field">
      <label class="label">Username or Email address</label>
      <div class="control">
        <input class="input" type="text" v-model="credential" autocomplete="current-username">
      </div>
    </div>
    <div class="field">
      <label class="label">Password</label>
      <div class="control">
        <input class="input" type="password" v-model="password" autocomplete="current-password">
      </div>
    </div>
    <p v-if="!success" class="help is-danger">{{ status }}</p> <!-- Display error message here -->
    <div class="field is-grouped">
      <div class="control">
        <button type="submit" class="button is-link">Log in</button>
      </div>
    </div>
  </form>
</template>
