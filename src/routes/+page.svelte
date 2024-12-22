<script lang='ts'>
  import clsx from 'clsx';

  const { data } = $props();
  const { user } = $derived(data);
  let selectedFiles = $state<File[]>([]);

  const upload = async () => {
    for (const file of selectedFiles) {
      const arrayBuffer = await file.arrayBuffer();
      const hashBuffer = await crypto.subtle.digest(
        'SHA-256',
        arrayBuffer,
      );
      const hashArray = Array.from(new Uint8Array(hashBuffer));
      const hashHex = hashArray
        .map(b => b.toString(16).padStart(2, '0'))
        .join('');
      const res = await fetch('/upload', {
        method: 'POST',
        body: JSON.stringify({
          name: file.name,
          size: file.size,
          hash: hashHex,
        }),
      });
      const data = await res.json();
      await fetch(data.url, {
        method: 'PUT',
        body: file,
      });
    }
  };
</script>

<div class='container is-max-tablet px-4'>
  <div class='image-upload-title'>
    <h1 class='title is-4'>Image Upload</h1>
    <p>5 MB max per file. 10 files max per request.</p>
  </div>
  <div>
    <div class='file is-boxed file-cta image-upload-height'>
      <label
        for='file-upload'
        class={clsx('container file-label is-responsive', { 'cursor-default': !user })}
      >
        <span class='image-upload-content'>
          {#if !user}
            <div class='has-text-centered'>
              You need to <a href='/signin'>sign in</a> before you can upload images
            </div>
          {:else if selectedFiles.length > 0}
            {#each selectedFiles as file}
              <div>
                <img
                  src={URL.createObjectURL(file)}
                  alt={file.name}
                  width={200}
                  height={200}
                />
                <span class='file-name is-small' title={file.name}>
                  {file.name}
                </span>
              </div>
            {/each}
          {:else}
            <div class='has-text-centered'>
              Drag &amp; drop files here ...
              <br />
              or
              <br />
              Copy &amp; paste screenshots here ...
            </div>
          {/if}
        </span>
      </label>
    </div>
    {#if user}
      <div class='is-flex'>
        <label class='file-label container select-files-label'>
          <input
            id='file-upload'
            class='file-input select-files-input'
            type='file'
            name='image'
            accept='.jpg, .jpeg, .png, .gif, .bmp'
            multiple
            onchange={(event: Event) => {
              const target = event.target as HTMLInputElement;
              selectedFiles = Array.from(target.files ?? []);
            }}
          />
          <span class='select-files'>Select images...</span>
        </label>
        {#if selectedFiles.length > 0}
          <div class='is-flex'>
            <button
              class='button is-danger select-button'
              type='button'
              onclick={() => (selectedFiles = [])}
            >
              Clear
            </button>
            <button
              class='button is-primary select-button'
              onclick={upload}
            >
              Upload
            </button>
          </div>
        {/if}
      </div>
    {/if}
  </div>
</div>

<style>
  .image-upload-title {
    border-bottom: 1px solid;
    margin-top: 1rem;
    margin-bottom: 1rem;
  }

  .image-upload-height {
    height: 100%;
    min-height: 20rem;
  }

  .image-upload-content {
    margin-top: auto;
    margin-bottom: auto;
    overflow: auto;
  }

  .select-files-label {
    border: 1px solid;
    border-color: var(--bulma-border);
  }

  .select-files {
    display: flex;
    align-items: center;
  }

  .select-files-input {
    cursor: pointer;
  }

  .cursor-default {
    cursor: default;
  }

  .select-button {
    border-radius: 0;
  }

  img {
    height: 200px;
    width: auto;
  }
</style>
