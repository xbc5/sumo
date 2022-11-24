package database_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/xbc5/sumo/internal/pkg/database"
	"github.com/xbc5/sumo/internal/pkg/database/model"
	t "github.com/xbc5/sumo/internal/pkg/mytest"
)

func expectPatternTags(fixture model.Pattern, result model.Pattern, tagLen int) {
	Expect(result.Name).To(Equal(fixture.Name))
	Expect(result.Description).To(Equal(fixture.Description))
	Expect(result.Pattern).To(Equal(fixture.Pattern))

	Expect(result.Tags).To(HaveLen(tagLen))
	for i, tag := range result.Tags {
		Expect(tag.Name).To(Equal(fixture.Tags[i].Name))
	}
}

var _ = Describe("AddPattern() and GetPatterns()", func() {
	Context("AddPattern() adds multiple records", func() {
		It("should not error", func() {
			db := t.OpenDB()
			fixtures := []model.Pattern{
				t.FakePattern(0, ".pattern0.", []string{"tag0", "tag1"}),
				t.FakePattern(1, ".pattern1.", []string{"tag1", "tag2", "tag3"}),
				t.FakePattern(2, ".pattern2.", []string{"tag1", "tag2", "tag3"}),
			}

			for _, f := range fixtures {
				addErr := database.AddPattern(db, f)
				Expect(addErr).ShouldNot(HaveOccurred())
			}

			database.Close(db)
		})
	})

	Context("AddPattern() sets one record; GetPatterns() retrieves all patterns", func() {
		It("should retrieve all patterns, preserving their relationships to tags", func() {
			db := t.OpenDB()
			pattern0 := t.FakePattern(0, ".pattern0.", []string{"tag0", "tag1"})

			database.AddPattern(db, pattern0)

			result, getErr := database.GetAllPatterns(db)
			Expect(getErr).ShouldNot(HaveOccurred())

			Expect(result).To(HaveLen(1))
			expectPatternTags(pattern0, result[0], 2)

			database.Close(db)
		})
	})

	Context("AddPattern() sets two records; GetPatterns() retrieves all patterns", func() {
		It("should retrieve all patterns, preserving their relationships to tags", func() {
			db := t.OpenDB()
			fixture0 := t.FakePattern(0, ".pattern0.", []string{"tag0", "tag1"})
			fixture1 := t.FakePattern(1, ".pattern1.", []string{"tag1", "tag2", "tag3"})

			database.AddPattern(db, fixture0)
			database.AddPattern(db, fixture1)

			result, getErr := database.GetAllPatterns(db)
			Expect(getErr).ShouldNot(HaveOccurred())

			Expect(result).To(HaveLen(2))

			expectPatternTags(fixture0, result[0], 2)
			expectPatternTags(fixture1, result[1], 3)

			database.Close(db)
		})
	})

	Context("AddPattern() adds existing record", func() {
		It("should not duplicate, and return only one record", func() {
			db := t.OpenDB()
			pattern0 := t.FakePattern(0, ".pattern0.", []string{"tag0", "tag1"})
			pattern1 := t.FakePattern(0, ".pattern0.", []string{"tag0", "tag1"})

			addErr := database.AddPattern(db, pattern0)
			Expect(addErr).ShouldNot(HaveOccurred())

			addErr = database.AddPattern(db, pattern1)
			Expect(addErr).ShouldNot(HaveOccurred())

			result, _ := database.GetAllPatterns(db)

			Expect(result).To(HaveLen(1))

			database.Close(db)
		})
	})

	Context("AddPattern() different tags", func() {
		It("should maintain associations to existing tags", func() {
			db := t.OpenDB()

			pattern0 := t.FakePattern(0, ".pattern0.", []string{"tag0"})
			pattern1 := t.FakePattern(0, ".pattern0.", []string{"tag1"})

			database.AddPattern(db, pattern0)
			database.AddPattern(db, pattern1)

			result, _ := database.GetAllPatterns(db)
			tags := result[0].Tags
			tag0 := tags[0].Name
			tag1 := tags[1].Name

			Expect(tags).To(HaveLen(2))
			for _, tag := range tags {
				Expect(tag.Name).To(BeElementOf([]string{"tag0", "tag1"}))
			}

			Expect(tag0).NotTo(Equal(tag1))

			database.Close(db)
		})
	})
})
