<script lang='ts'>
  import type { FileWithPreview } from '$lib/types';
  import { calculateFileHash } from '$lib/utils/file';
  import { onDestroy, onMount } from 'svelte';

  const { user, files, onFilesSelected, onRemove, onClear } = $props<{
    user: any;
    files: FileWithPreview[];
    onFilesSelected: (files: FileWithPreview[]) => void;
    onRemove: (index: number) => void;
    onClear: () => void;
  }>();

  let dragCounter = $state(0);

  function handleGlobalPaste(e: ClipboardEvent) {
    if (!user)
      return;
    const pastedFiles = Array.from(e.clipboardData?.files ?? [])
      .filter(file => file.type.startsWith('image/'));
    if (pastedFiles.length > 0) {
      getUniqueFiles(pastedFiles).then((uniqueFiles) => {
        if (uniqueFiles.length > 0) {
          onFilesSelected(uniqueFiles);
        }
      });
    }
  }

  onMount(() => {
    if (import.meta.env.SSR)
      return;
    document.addEventListener('paste', handleGlobalPaste);
  });

  onDestroy(() => {
    if (import.meta.env.SSR)
      return;
    document.removeEventListener('paste', handleGlobalPaste);
  });

  function handleDragEnter(e: DragEvent) {
    e.preventDefault();
    dragCounter++;
  }

  function handleDragLeave(e: DragEvent) {
    e.preventDefault();
    dragCounter--;
  }

  async function createFileWithPreview(file: File): Promise<FileWithPreview> {
    if (!file || !(file instanceof File)) {
      throw new Error('Invalid file');
    }
    const hash = await calculateFileHash(file);
    const preview = URL.createObjectURL(file);

    return { file, preview, hash };
  }

  async function getUniqueFiles(filesArray: File[]): Promise<FileWithPreview[]> {
    const seenHashes = new Set<string>(files.map((existingFile: FileWithPreview) => existingFile.hash));

    const filesWithHash = await Promise.all(filesArray.map(async (file) => {
      if (file instanceof File) {
        return createFileWithPreview(file);
      }
      else {
        console.error('Invalid file:', file);
        return null;
      }
    }));

    return filesWithHash.filter((newFile): newFile is FileWithPreview => {
      if (newFile === null)
        return false;

      if (seenHashes.has(newFile.hash)) {
        return false;
      }

      seenHashes.add(newFile.hash);
      return true;
    });
  }

  async function handleDrop(e: DragEvent) {
    e.preventDefault();
    dragCounter = 0;
    if (!user)
      return;

    const droppedFiles = Array.from(e.dataTransfer?.files ?? [])
      .filter(file => file.type.startsWith('image/'));
    if (droppedFiles.length > 0) {
      const uniqueFiles = await getUniqueFiles(droppedFiles);
      if (uniqueFiles.length > 0) {
        onFilesSelected(uniqueFiles);
      }
    }
  }

  function removeFile(index: number) {
    onRemove(index);
  }

  async function upload() {
    for (const { file } of files) {
      const hash = await calculateFileHash(file);
      const res = await fetch('?/upload', {
        method: 'POST',
        body: JSON.stringify({
          name: file.name,
          size: file.size,
          hash,
        }),
      });
      const result = await res.json();
      if (result.type !== 'success') {
        console.error('Upload failed:', result.error);
        continue;
      }
      const { url } = result.data;
      await fetch(url, {
        method: 'POST',
        body: file,
      });
    }
    window.location.reload();
  }
</script>

<div>
  <div class='file is-boxed file-cta image-upload-height'>
    <div
      class='container file-label is-responsive'
      class:cursor-default={!user}
      class:is-dragover={dragCounter > 0}
      role='button'
      tabindex='0'
      ondragenter={handleDragEnter}
      ondragleave={handleDragLeave}
      ondragover={e => e.preventDefault()}
      ondrop={handleDrop}
    >
      <span class='image-upload-content'>
        {#if !user}
          <div class='has-text-centered'>
            You need to <a href='/signin'>sign in</a> before you can upload images
          </div>
        {:else if files.length > 0}
          {#each files as { preview, file }, i}
            <div class='preview-container'>
              <button
                type='button'
                class='delete-button'
                onclick={(e) => {
                  e.stopPropagation();
                  removeFile(i);
                }}
                aria-label='Delete image'
              >
                ×
              </button>
              <img
                src={preview}
                alt={file.name}
                class='preview-image'
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
    </div>
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
          onchange={async (e: Event) => {
            const target = e.target as HTMLInputElement;
            const filesArray = Array.from(target.files ?? []);
            const uniqueFiles = await getUniqueFiles(filesArray);
            if (uniqueFiles.length > 0) {
              onFilesSelected(uniqueFiles);
            }
          }}
        />
        <span class='select-files'>Select images...</span>
      </label>

      {#if files.length > 0}
        <div class='is-flex'>
          <button
            class='button is-danger select-button'
            type='button'
            onclick={onClear}
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

<style>
  .preview-container {
    position: relative;
    display: inline-block;
    margin: 0.5rem;
    border-radius: 0.25rem;
    overflow: hidden;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .delete-button {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    width: 1.5rem;
    height: 1.5rem;
    border-radius: 50%;
    background: rgba(0, 0, 0, 0.5);
    color: white;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1rem;
    transition: background-color 0.2s;
  }

  .delete-button:hover {
    background: rgba(0, 0, 0, 0.7);
  }

  .preview-image {
    width: 200px;
    height: 200px;
    object-fit: cover;
    display: block;
  }

  .image-upload-height {
    height: 100%;
    min-height: 20rem;
    border: 2px dashed var(--bulma-border);
    border-radius: 0.5rem;
    transition: all 0.2s ease;
  }

  .image-upload-height:hover,
  .is-dragover {
    border-color: var(--bulma-primary);
    background-color: rgba(var(--bulma-primary-rgb), 0.05);
  }

  .image-upload-content {
    margin-block: auto;
    overflow: auto;
    padding: 1rem;
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    justify-content: center;
  }

  .select-files-label {
    border: 1px solid var(--bulma-border);
    border-radius: 0.25rem;
    overflow: hidden;
  }

  .select-files {
    display: flex;
    align-items: center;
    padding: 0.5rem 1rem;
  }

  .select-files-input {
    cursor: pointer;
  }

  .cursor-default {
    cursor: default;
  }

  .select-button {
    margin-left: 0.5rem;
  }
</style>
