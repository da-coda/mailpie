<template>
  <v-layout>
  </v-layout>
</template>

<script>
export default {
  name: "Home",
  data: () => {
    return {
      receivedMails: []
    }
  },
  computed: {
    mails() {
      return this.receivedMails
    }
  }
}
</script>