import { marked } from 'marked'
import DOMPurify from 'dompurify'

DOMPurify.addHook('afterSanitizeAttributes', function (node) {
  // set all elements owning target to target=_blank
  if ('target' in node) {
    node.setAttribute('target', '_blank')
    node.setAttribute('rel', 'noopener')
  }
})

export default function(markdown) {
  return DOMPurify.sanitize(marked.parse(markdown))
}
