<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import HelloWorld from './components/HelloWorld.vue'
import {onMounted} from "vue";
import { useApiClientStore } from './stores/apiClient.ts';

const apiClientStore = useApiClientStore();

onMounted(() => {
  apiClientStore.initFromStorage();
})
</script>

<template>
  <header>
    <img alt="Vue logo" class="logo" src="@/assets/logo.svg" width="125" height="125" />

    <div class="wrapper">
      <nav>
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/login">About</RouterLink>
      </nav>
    </div>
  </header>

  <RouterView />
</template>
