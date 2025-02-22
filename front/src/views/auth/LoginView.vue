<script setup lang="ts">
  import { useApiClientStore } from "@/stores/apiClient.ts";
  import LoginForm from "@/components/auth/LoginForm.vue";
  import {RouterLink} from "vue-router";
  const apiClientStore = useApiClientStore();

  let authenticated = apiClientStore.isAuthenticated();
</script>

<template>
  <main>
    <section v-if="!authenticated" class="section">
      <h1 class="title">Log in</h1>
      <LoginForm/>
    </section>
    <section v-if="authenticated" class="section">
      <h1 class="title">Already logged in!</h1>
      <RouterLink class="button is-link" to="/">Back to home</RouterLink>
    </section>
  </main>
</template>
