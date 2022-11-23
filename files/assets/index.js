const prefix = 'switchweb'

const onError = async (msg) => {
  document.getElementById('alert').style.display = null
  document.getElementById('alertMsg').innerHTML = msg
}

const onAct = async (ev) => {
  if (!ev.currentTarget) {
    return
  }
  const target = ev.currentTarget
  const id = target.dataset.id

  if (ev.altKey) {
    target.style.display = 'none'
    localStorage.setItem(`${prefix}/button/${id}/hide`, true)
    ev.stopPropagation()
    return
  }

  target.style.backgroundColor = '#ccc'
  try {
    const res = await fetch(`act/${id}`, {
      method: 'POST',
      cache: 'no-cache',
    })
    if (!res.ok) {
      const body = await res.text()
      onError(`fail: ${body}`)
    }
  } catch (err) {
    onError(`error: ${err}`)
  }
  target.style.backgroundColor = null
}

const onDevice = async (ev) => {
  if (!ev.currentTarget) {
    return
  }
  const id = ev.currentTarget.dataset.id

  if (ev.altKey) {
    ev.currentTarget.style.display = 'none'
    localStorage.setItem(`${prefix}/device/${id}/hide`, true)
    return
  }
}

Array.from(document.getElementsByClassName('action')).forEach((e) => {
  const id = e.dataset.id
  e.addEventListener('click', onAct)
  const hide = localStorage.getItem(`${prefix}/button/${id}/hide`)
  if (hide) {
    e.style.display = 'none'
  }
})

Array.from(document.getElementsByClassName('device')).forEach((e) => {
  const id = e.dataset.id
  e.addEventListener('click', onDevice)
  const hide = localStorage.getItem(`${prefix}/device/${id}/hide`)
  if (hide) {
    e.style.display = 'none'
  }
})

document.getElementById('restore').addEventListener('click', (ev) => {
  if (!ev.target) {
    return
  }
  const id = ev.target.dataset?.id
  if (!id) {
    return
  }
  localStorage.removeItem(`${prefix}/device/${id}/hide`)
  Array.from(document.getElementsByClassName('device')).
    filter((e) => e.dataset.id == id).
    forEach((e) => {
      e.style.display = null
      Array.from(e.getElementsByClassName('action')).forEach((e) => {
        const id = e.dataset.id
        localStorage.removeItem(`${prefix}/button/${id}/hide`)
        e.style.display = null
      })
    })
})

const closeAlert = () => {
  document.getElementById('alert').style.display = 'none'
  document.getElementById('alertMsg').innerHTML = null
}
document.getElementById('alert').addEventListener('click', closeAlert)
closeAlert()
