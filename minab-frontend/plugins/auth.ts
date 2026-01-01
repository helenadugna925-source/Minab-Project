export default defineNuxtPlugin((nuxtApp) => {
  const cookie = useCookie('token')

  // The hook gives us 'token', which is a Ref we can update
  nuxtApp.hook('apollo:auth', ({ token }) => {
    // If the cookie has a value, send the Bearer token.
    // If NOT, setting it to null tells Hasura to use the 'anonymous' role.
    if (cookie.value && cookie.value !== 'undefined') {
      token.value = `Bearer ${cookie.value}`
    } else {
      token.value = null
    }
  })
})