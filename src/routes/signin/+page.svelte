<form method='post'>
  <div class='field'>
    <label for='email'>Email:</label>
    <input
      class='input'
      id='email'
      name='email'
      type='email'
      required
      placeholder={import.meta.env.VITE_EMAIL_PLACEHOLDER || ''} />
  </div>
  <div class='field'>
    <label for='password'>Password:</label>
    <input
      class='input'
      id='password'
      name='password'
      type='password'
      required
      placeholder={import.meta.env.VITE_PASSWORD_PLACEHOLDER || ''} />
  </div>
  <button class='button' type='submit'> Sign in </button>
</form>
