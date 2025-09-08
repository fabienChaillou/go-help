## Summary golang concpets

### keys word

Go language **imperatif et non-fonctionnel**

atomicity:
    fait référence à l'exécution d'une opération ou d'un ensemble d'opérations de manière **indivisible**
        * **transactions en base de données**
        * **concurrences**

rollback:

channel:
    Les **channels** en Go sont des **tuyaux de communication** entre les goroutines. Ils permettent de **transmettre des données de manière sûre** et **synchronisée**, sans avoir à utiliser de mutex dans la plupart des cas.

buffer:
    En Go, le **buffer** fait généralement référence à un espace mémoire temporaire utilisé pour **stocker des données avant traitement ou transmission**. Le concept de buffer est souvent utilisé avec :

        1. **Les channels bufferisés** (concurrence)
        2. **Les buffers d'entrée/sortie** (fichiers, réseaux, etc.)

Interopérabilité:
    L’interopérabilité avec C en Go se fait via un mécanisme appelé **`cgo`**. Il te permet d’appeler du code écrit en **C** directement depuis du code Go.

Inference de type
     est un mécanisme par lequel le compilateur déduit automatiquement le type d'une variable à partir de sa valeur initiale, sans que tu aies besoin de le spécifier explicitement.

**`sync.WaitGroup`** `ync.WaitGroup`
    est un outil de synchronisation qui permet d’attendre que plusieurs goroutines aient terminé leur travail.

**mutex** (abréviation de *mutual exclusion*) `var mutex sync.Mutex`
    en Go est un outil de synchronisation qui permet de **protéger l'accès à une ressource partagée** (comme une variable ou une structure) contre les accès concurrents par plusieurs goroutines. (**data race)

**`os.Pipe()`** permet de créer un **pipe anonyme**, 
    c’est-à-dire un mécanisme de communication **unidirectionnelle** entre deux parties d’un programme (souvent entre deux goroutines ou entre un parent et un sous-processus).

`defer` permet de **différer** l'exécution d'une fonction jusqu'à la fin de la fonction englobante,

`select`
    une instruction utilisée pour écouter plusieurs canaux (`chan`) en même temps. Elle permet de gérer la concurrence de manière élégante en attendant qu’un ou plusieurs canaux soient prêts pour la communication (lecture ou écriture). Elle est souvent utilisée avec les goroutines pour synchroniser les opérations concurrentes.

**package `context`** 
    utilisé pour **transporter des délais d'expiration, des annulations et d'autres valeurs spécifiques à la requête** à travers les appels de fonctions et les goroutines. Il est essentiel pour gérer les **opérations asynchrones**, les **requêtes HTTP**, et les **bases de données** de manière contrôlée.

**worker pool** — 
    un modèle très courant et puissant en Go pour **traiter plusieurs tâches en parallèle**, tout en contrôlant **le nombre de goroutines**.

La **récursivité** (ou factoriel)
    en Go (ou Golang), comme dans d'autres langages, est une technique où une fonction s'appelle **elle-même** pour résoudre un problème. C'est souvent utilisé pour résoudre des problèmes décomposables en sous-problèmes similaires (ex. : calcul du facteur d'une valeur, parcours d'arbre, etc.).

les driver db
[sql-driver](https://go.dev/wiki/SQLDrivers)

[pgx](https://github.com/jackc/pgx)
