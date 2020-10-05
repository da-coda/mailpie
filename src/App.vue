<template>
  <v-app>
    <v-app-bar app>
      <h1>Mailmon Go</h1>
    </v-app-bar>
    <v-main>
      <v-container fluid fill-height>
        <router-view></router-view>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>

export default {
  name: "App",
  components: {

  }
};
</script>

<style></style>