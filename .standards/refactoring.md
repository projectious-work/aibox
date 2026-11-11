# Refactoring standard

## Purpose and scope

This standard defines how projectious.work software changes structure
without changing behavior: breaking up large files, extracting libraries
from monoliths, and paying down technical debt. It applies to every
refactoring larger than a local rename.

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

## Principles

- Refactoring changes structure, never behavior. Behavior changes are
  separate work with their own review.
- The build and the test suite stay green at every commit.
- A refactoring has a named target, not a general urge to "clean up".

## Procedure

1. **Name the target.** State the *Code Smells* to remove (for example Long
   Method, Feature Envy, Shotgun Surgery, Divergent Change) and the
   refactorings from the *Refactoring Catalog according to Fowler* to apply.
2. **Know why the code is as it is** before removing a check, boundary, or
   oddity (*Chesterton's Fence*).
3. **Secure behavior first.** If existing tests do not cover the code at its
   public boundary, add characterization tests before changing structure.
4. **Move in small steps.** Each step is one catalog refactoring followed by
   the affected test run.
5. **Larger changes use the *Mikado Method*.** Record the Mikado
   graph in the pull request or development note. Never leave the build broken
   while digging.
6. **Respect the target architecture.** Refactorings move code towards the
   [software architecture and code design standard](software-architecture-and-code-design.md),
   never across its boundaries in the wrong direction. Extractions produce
   *Deep Modules*, not shallow pass-through layers.
7. **Keep scope.** Clean up code you touch anyway (*Boy Scout Rule*); larger
   findings become explicit follow-up work.

## Commits and review

- Refactoring happens on `chore/` branches; its commits use the `refactor:` type
  ([Conventional Commits](git-branching-and-releases.md#commit-messages))
  and contain no behavior change.
- Generated output, golden fixtures, and public contracts are unchanged by a
  refactoring. A diff in them is a behavior change and stops the refactoring.
- The pull request lists smells removed, refactorings applied, tests run, and
  deferred findings.

## Agent briefing

> Name the *Code Smells*, learn why code exists (*Chesterton's Fence*), secure
> behavior with characterization tests, apply *Refactoring Catalog* steps one
> at a time with green tests at every commit, use the *Mikado Method* for large
> changes, extract *Deep Modules*, limit scope by the *Boy Scout Rule*, and
> keep behavior changes out of `refactor:` commits.
