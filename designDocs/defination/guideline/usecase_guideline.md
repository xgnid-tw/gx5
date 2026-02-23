# Use Case Definition Guidelines

## Purpose

This guideline provides a framework for creating consistent use case definition documents in software development projects.

## Target Audience

- System requirements analysts
- Business analysts
- Project members

---

## 1. Use Case Document Creation Policy

### 1.1 Scope of Each Document

- **One function, one use case:** Create one use case definition document per business function.
- **Consistent granularity:** Describe similar functions at the same level of detail.
- **Independence:** Each use case should be self-contained without relying on other use cases.

### 1.2 Document Separation Policy

To improve maintainability and readability, project documents are separated as follows:

#### Use Case Definition Document (UC-XXX)

- **Scope:** Use case overview, actors, pre/post-conditions, business rules
- **Purpose:** Define system requirements
- **Audience:** Developers, testers, system designers

#### Business Flow Definition Document (BF-XXX)

- **Scope:** Detailed business flows, data requirements, screen/API specifications
- **Purpose:** Define detailed business processes
- **Audience:** Business stakeholders, UI designers, detailed design engineers

### 1.3 Naming Conventions

#### Use Case Names

- **Verb + Object format:** e.g., "Register Invoice," "Submit Payment Request," "Update Company Information"
- **Business-oriented naming:** Name from the user's business perspective, not from internal system processing.
- **No abbreviations:** Use "User Registration" instead of "UC Registration."

#### File Names

- **Use Case Definition:** `UC-XXX_[Use Case Name].md`
- **Business Flow Definition:** `BF-XXX-Y_[Use Case Name]_Business_Flow_[Scenario].md`
- **Examples:**
  - `UC-001_Invoice_Registration.md` → `defination/usecases/UC-001_Invoice_Registration.md`
  - `BF-001-1_Invoice_Registration_Business_Flow_Normal.md` → `defination/flow/BF-001-1_Invoice_Registration_Business_Flow_Normal.md`
  - `BF-001-2_Invoice_Registration_Business_Flow_Input_Error.md` → `defination/flow/BF-001-2_Invoice_Registration_Business_Flow_Input_Error.md`

#### ID System

- **Use Case ID:** `UC-XXX` (3-digit sequential number)
- **Business Flow ID:** `BF-XXX-Y` (XXX: use case number, Y: sequential number within the same use case)
- **Numbering Rules:**
  - Assign use case numbers in ranges by functional category.
  - Business flow sequential numbers start at 1 for the normal flow, with 2 and above for alternative/exception flows.
  - Refer to the [Use Case Index](../usecases/Use_Case_Index.md) for detailed number assignments.

## 2. Section-by-Section Writing Guidelines

### 2.1 Document Metadata

- **Required fields:** Fill in all fields.
- **Status management:** Manage in the order Draft → Review → Approved.
- **Version control:** Increment the integer part for major changes and the decimal part for minor changes.

### 2.2 Use Case Overview

#### Purpose

- **State the business value:** Describe why this use case is necessary.
- **Clarify deliverables:** Describe what is achieved upon completion of the use case.

#### Summary

- **1–2 sentences, concise:** Detailed enough for a third party to understand.
- **Clear subject and predicate:** Clearly state "who does what."

#### Scope

- **Define the covered range:** Explicitly state the boundaries of functions handled by this use case.
- **State exclusions:** Clearly document what is out of scope to avoid confusion.

### 2.3 Actor Information

#### Actor Definitions

Use actor names defined in the [Use Case Index](../usecases/Use_Case_Index.md).

#### Primary Actor

- **Business executor:** The person or system that actually executes the use case.
- **Only one:** There should be only one primary actor per use case.

#### Secondary Actor

- **Supporting role:** Person or system that supports the primary actor.
- **Include as needed:** If none exist, state "None."

#### System Actor

- **External systems:** External systems or APIs that interact with this use case.
- **Clarify implementation scope:** Constraints for specific implementation phases (e.g., MVP) may be noted if needed (generic descriptions are recommended by default).

### 2.4 Pre-conditions and Post-conditions

#### Pre-conditions

- **State prerequisite states:** Required states before the use case begins.
- **Confirm data existence:** Explicitly note required master data, etc.
- **Authentication and authorization:** Specify user login status and permission levels.

#### Post-conditions

- **Describe success and failure separately:** Document the resulting state for both outcomes.
- **Specify data changes:** Document which data changes and how.
- **Specify state transitions:** Document state changes of the system and business objects.

### 2.5 Business Flows

#### Summary Flow

- **Basic processing sequence:** Briefly describe the main processing flow of the use case.
- **Step numbering:** Number each step to clarify order.
- **Identify actors:** Indicate which actor acts at each step (at a summary level).
- **Readability focus:** Write at a level understandable by non-technical readers.

#### Detailed Business Flows

- **Refer to separate documents:** Detailed flows (basic, alternative, exception) are documented in Business Flow Definition Documents (BF-XXX).
- **Cross-references:** Include clear reference links to the Business Flow Definition Documents.
- **When no business usage is defined:** If no specific business usage has been identified, state: "At this time, no specific business usage calling this function has been identified; therefore, a detailed business flow definition is not provided." Do not create a Business Flow Definition Document in this case.
- **Benefits of document separation:**
  - Use Case Definition: Focuses on system requirements.
  - Business Flow Definition: Focuses on detailed processing steps.

### 2.6 Business Rules

#### Business Rule Description

- **Unique identification:** Assign a unique ID in the format BR-XXX.
- **Include rule name:** Use the format "BR-XXX: Rule Name" combining ID and rule name.
- **Concise description:** Describe rule content in 1–2 sentences.
- **Judgment criteria:** Describe objectively verifiable criteria.
- **Exception conditions:** Document exception conditions where the rule does not apply.

#### Business Rule Numbering System

Assign unique IDs in the BR-XXX format. Verify no duplicate IDs exist across all use case definitions.

#### Business Rule Cross-References

- **Mandatory cross-referencing:** Cross-references between Use Case Definitions and the Use Case Index are required.
- **Synchronized updates:** When adding or changing business rules in a Use Case Definition, ensure consistency with other use cases.
- **Consistency assurance:** Periodically verify consistency of business rule IDs across all Use Case Definitions.

### 2.7 Related Use Cases

#### Include Relationships

- **Mandatory includes:** The following infrastructure functions must be included based on conditions:
  - UC-902: Logging — Mandatory include only for operations involving data changes (create, update, delete, submit, etc.); not required for simple read operations.
  - UC-903: Error Handling — Mandatory include for all use cases.
  - Note: UC-901 (User Authentication) is included only in use cases that assume the user is already authenticated.
- **Business-level includes:** Other use cases required for executing this use case.
- **Relationship description:** Briefly explain why the relationship is necessary.

#### Extend Relationships

- **Conditional extensions:** Extension use cases executed under specific conditions.
- **Extension conditions:** Clearly state the conditions under which the extension occurs.

#### Generalization Relationships

- **Abstraction of commonalities:** Used when abstracting common parts across multiple use cases.
- **Inheritance:** Document the relationship with the parent use case.

### 2.8 Supplementary Information

#### Expected Usage Frequency

- **Usage frequency:** Specify concrete figures (daily, monthly, etc.).
- **When no usage is expected:** If there is no specific usage expectation, state: "No usage expected at this time."
- **Peak hours:** Explain peak hours and their reasons.
- **Impact on system design:** Consider the impact on performance requirements and resource planning.

#### Operations and Maintenance Requirements

- **Data retention period:** Retention period for related data.
- **Maintenance notes:** Key points for operations and maintenance.
- **System integration:** Operational requirements related to integration with external systems.
- **When no requirements exist:** If there are no specific operations/maintenance requirements, state: "None."

#### Other Notes

- **Feature constraints:** Clearly describe features not provided in the current system.
- **Future extensions:** Document features planned for future addition.
- **Technical constraints:** Technical constraints or special implementation requirements.

## 3. Project-Specific Considerations

### 3.1 External System Integration

#### Integration Pattern Classification

- **API Integration Patterns**
  - Synchronous API: Functions requiring real-time processing (e.g., payment submission/cancellation)
  - Asynchronous API: Functions that can be handled by batch processing (e.g., data ingestion)
- **Screen Transition Patterns**
  - Seamless transition: Integrated authentication via SSO
  - Separate screen transition: Transition to an external service screen
- **Notification-Based Patterns**
  - Webhook: Real-time notification reception
  - Polling: Periodic check

#### Integration Specification Documentation

- **API specification references:** Include references to the API specification documents for integrations.
- **Error handling:** Document processing methods when external systems fail.
- **Data synchronization:** Document data synchronization timing and frequency.
- **Authentication methods:** Document the authentication method for each external system.

### 3.2 Permission Management Standardization

#### Permission Management Patterns

- **Role-based feature restriction:** Separation of submission and approval roles.
- **Organization-level data access control:** Access limited to own organization's data.
- **Self-approval prohibition rule:** Users cannot approve their own submissions.

#### Permission Requirement Documentation

- **Permissions in pre-conditions:** Specify required permission levels.
- **Access control rules:** Describe data access restrictions concretely.
- **Separation of operational permissions:** Specify separation of read, write, and approval permissions.

### 3.3 Non-Functional Requirements in Supplementary Information

Non-functional requirements are not documented as a standalone section but are included within the Supplementary Information section as needed.

#### Performance Examples

- **Expected usage frequency:** e.g., Average 50 transactions per organization per month, up to 100 at peak.
- **Response time:** Document in operations/maintenance requirements as needed.
- **When no requirements exist:** If there are no specific performance requirements, state: "None."

#### Security Examples

- **Data retention period:** Document in operations/maintenance requirements (e.g., invoice data retained for 10 years).
- **Permission requirements:** Document in the pre-conditions section.
- **When no requirements exist:** If there are no specific security requirements, state: "None."

#### Other Notes Examples

- **Technical constraints:** System-specific constraints.
- **Future extensions:** Document features planned for future addition.
- **Relationship with existing features:** Document the relationship and distinction between existing similar features in the current system.
- **Provisional content:** Document content that is provisional due to unconfirmed details, along with plans for future refinement.
- **Implementation phase constraints:** Constraints for specific implementation phases may be noted if needed (generic descriptions are recommended by default).

#### Example: Relationship with Existing Features

```markdown
- **Relationship with existing features:** The account transfer approval function has existing implementations in System A and System B, but their detailed specifications have not been confirmed. The content in this document may contain inaccuracies.
```

#### Example: Provisional Content

```markdown
- **Provisional content:** The content related to Feature X approval in this document is provisional and is expected to be refined or revised through further detailed investigation and analysis.
- **Items to be detailed in the future:**
  - Feature X-specific approval requirements (e.g., credit checks, third-party integration)
  - Consistency verification with existing System A approval functions
  - Specific screen specifications and operational flows in System B
```

### 3.4 Status Transitions

- **Status management:** Document state management (e.g., invoice status).
- **Transition conditions:** Document the conditions under which status changes occur.
- **Consistency with transition diagrams:** Ensure consistency with separately created status transition diagrams.

## 4. Quality Checkpoints

### 4.1 Completeness Check

- [ ] Are all required sections documented?
- [ ] Are pre-conditions and post-conditions clearly defined?
- [ ] Are major alternative and exception flows covered?
- [ ] Are relationships with related use cases documented?

### 4.2 Consistency Check

- [ ] Is terminology consistent with other use case definitions?
- [ ] Are actor names unified?
- [ ] Are data and field names unified?
- [ ] Are status names unified?
- [ ] Are business rule numbers free of duplicates with other use cases?

#### Consistency Check with Use Case Index

- [ ] Do the use case ID and name match the Use Case Index?
- [ ] Does the primary actor match the Use Case Index entry?
- [ ] Is the use case summary consistent with the Use Case Index summary?
- [ ] Do actor names used match the actor definitions in the Use Case Index?
- [ ] Is the Use Case Index updated synchronously when creating, modifying, or deleting use cases?

#### Related Document Consistency Check

- [ ] Consistency with mandatory update documents (see Section 6.2.3 for details)
- [ ] Consistency of cross-references
- [ ] Recording of update history

### 4.3 Feasibility Check

- [ ] Is the content technically feasible?
- [ ] Is it achievable within the project scope?
- [ ] Is the external system integration approach appropriate?
- [ ] Are non-functional requirements realistic?

### 4.4 Testability Check

- [ ] Are requirements described in a testable manner?
- [ ] Are test perspectives adequately covered?
- [ ] Are test patterns for normal, abnormal, and boundary cases anticipated?

## 5. Review and Approval Process

### 5.1 Review Flow

1. **Author self-check:** Verify the quality checkpoints above.
2. **Team review:** Peer review within the same team.
3. **Stakeholder review:** Review by business stakeholders, external system owners, and other relevant parties.
4. **Final approval:** Final approval by the project manager.

### 5.2 Review Perspectives

- **Business perspective:** Is the business flow appropriate?
- **Technical perspective:** Is it technically feasible?
- **Project perspective:** Does it fit within the project scope?
- **Operational perspective:** Are there any issues from an operations/maintenance standpoint?
- **Consistency perspective:** Is consistency with related documents ensured? Are necessary synchronized updates complete?

## 6. Management Methods

### 6.1 File Management

- **Storage location:** `designDocs/defination/usecases/` folder
- **Naming conventions:** Follow the naming conventions described above.
- **Version control:** Use Git for version management.

### 6.2 Related Document Consistency Management

When creating, modifying, or deleting Use Case Definitions, the following management process is mandatory to maintain consistency with related documents.

#### 6.2.1 Related Document Classification

Related documents are classified into three priority levels based on impact and modification frequency, ensuring consistency in stages.

##### Priority 1: Mandatory Update Documents (Must be updated simultaneously)

| Document Name | Storage Location | Update Timing | Consistency Check Items |
|---|---|---|---|
| **Use_Case_Index.md** | `defination/usecases/` | Simultaneously with use case creation/modification/deletion | Match of ID, name, primary actor, summary |
| **Business Flow Definition (BF-XXX)** | `defination/flow/` | When summary flow changes | Match of flow numbers and use case references |

##### Priority 2: Important Update Documents (Update within 1 week)

| Document Name | Storage Location | Update Timing | Impact Check Items |
|---|---|---|---|
| **Functional_Requirements_Index.md** | `defination/functions/` | When use cases related to system functions change | Consistency of related business flows and API functions |

##### Priority 3: Impact Confirmation Documents (Confirm at next review)

If related design documents or specifications exist, conduct impact confirmation as needed.

#### 6.2.2 Consistency Management Process

##### Required Steps When Adding a New Use Case

1. **Create Use Case Definition**
   - Create the use case definition following the template.
   - Verify no duplicate use case IDs or business rule IDs.

2. **Simultaneously Update Priority 1 Documents**
   - Add an entry to the [Use Case Index](../usecases/Use_Case_Index.md).
   - Create a BF-XXX file in `defination/flow/` if a detailed business flow is needed.

3. **Check and Update Priority 2 Documents**
   - If system functions are involved: Consider adding to the [Functional Requirements Index](../functions/Functional_Requirements_Index.md).

4. **Check Priority 3 Documents**
   - Confirm impact on related design documents and specifications (as needed).

##### Required Steps When Modifying an Existing Use Case

1. **Update Use Case Definition**
   - Update the use case definition according to the changes.
   - Clarify the scope of impact (actors, business rules, flows).

2. **Synchronize Priority 1 Documents**
   - Synchronize the corresponding entry in the Use Case Index (name, primary actor, summary).
   - Update the corresponding BF-XXX file when the summary flow changes.

3. **Assess Impact on Priority 2 and 3 Documents**
   - Determine whether related documents need updating based on the scope of impact.
   - Update applicable documents if necessary.

##### Required Steps When Deleting a Use Case

1. **Delete Use Case Definition**
   - Delete the use case definition file.
   - Record the reason for deletion in the change history.

2. **Synchronize Deletion in Priority 1 Documents**
   - Remove the corresponding entry from the Use Case Index.
   - Consider deleting the corresponding BF-XXX file.

3. **Remove References from All Priority Documents**
   - Remove references to the deleted use case from the Functional Requirements Index.
   - Remove references from related design documents and specifications (as needed).
   - Record the scope of deletion impact in the update history.

#### 6.2.3 Consistency Check Methods

##### Periodic Consistency Check (Monthly)

Use the following checklist to verify consistency:

**Consistency Check with Use Case Index**

- [ ] Are all Use Case Definitions registered in the Use Case Index?
- [ ] Do use case IDs, names, primary actors, and summaries match?
- [ ] Are actor definitions unified?
- [ ] Are deleted use cases removed from the Use Case Index?

**Consistency Check with Business Flow Definitions**

- [ ] Is there consistency between summary flows and detailed business flows?
- [ ] Are BF-XXX numbers correctly assigned?
- [ ] Are reference links correctly set?

##### Change-Time Consistency Check (At each change)

As part of the change management process, verify the following:

**Impact Scope Confirmation**

- [ ] Has the impact on other use cases that reference the changed use case been confirmed?
- [ ] Have other use cases using the changed actor been confirmed?
- [ ] Have other use cases applying the changed business rule been confirmed?

**Synchronized Update Confirmation**

- [ ] Have mandatory update documents (Priority 1) been updated simultaneously?
- [ ] Has an update schedule been set for important update documents (Priority 2)?
- [ ] Has a confirmation schedule been recorded for impact confirmation documents (Priority 3)?

### 6.3 Change Management

- **Change requests:** Submit change requests via the change management ledger.
- **Impact assessment:** Evaluate the impact on related use cases and associated documents.
- **Approval process:** Conduct an approval process appropriate to the change content and scope of impact.
- **Consistency assurance:** Update related documents according to the consistency management process in Section 6.2.

### 6.4 Traceability

- **Requirements tracing:** Ensure traceability from business requirements to use cases.
- **Design tracing:** Ensure traceability from use cases to detailed designs.
- **Test tracing:** Ensure traceability from use cases to test cases.

---

**Revision History**
| Version | Date | Author | Description |
|---------|------|--------|-------------|
| 1.0 | YYYY/MM/DD | Requirements Team | Initial version — generalized from project-specific guideline |
| 1.1 | 2026/02/22 | — | Update directory paths to match project structure |
