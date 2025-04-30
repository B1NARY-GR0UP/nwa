(defun hello-world ()
  (format t "Hello, World!~%"))

;; Define a class-like structure
(defstruct greeter
  (message "Hello, Lisp!"))

;; Call the function
(hello-world)