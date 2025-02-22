<script setup lang="ts">
  import {useApiClientStore} from "@/stores/apiClient.ts";
  import router from "@/router";
  import {RouterLink} from "vue-router";

  const apiClientStore = useApiClientStore();
  const authenticated = apiClientStore.isAuthenticated();

  async function handleLogout() {
    try {
      await apiClientStore.logout();
      await router.push("/");
    } catch (e: any) {
      console.error(e.message)
    }
  }
</script>

<template>
  <main>
    <section v-if="authenticated" class="section">
      <h1 class="title">Log out</h1>
      <button class="button is-link" @click="handleLogout">Log me out</button>
    </section>
    <section v-if="!authenticated" class="section">
      <h1 class="title">You're not logged in</h1>
      <RouterLink class="button is-link" to="/login">Log in</RouterLink>
    </section>
  </main>
</template>
