package libxml2

import "testing"

// Tests for DOM Level 3

func TestDocumentAttributes(t *testing.T) {
	doc := CreateDocument()
	defer doc.Free()
	if doc.Encoding() != "" {
		t.Errorf("Encoding should be empty string at first, got '%s'", doc.Encoding())
	}

	if doc.Version() != "1.0" {
		t.Errorf("Version should be 1.0 by default, got '%s'", doc.Version())
	}

	if doc.Standalone() != -1 {
		t.Errorf("Standalone should be -1 by default, got '%d'", doc.Standalone())
	}

	for _, enc := range []string{"utf-8", "euc-jp", "sjis", "iso-8859-1"} {
		doc.SetEncoding(enc)
		if doc.Encoding() != enc {
			t.Errorf("Expected encoding '%s', got '%s'", enc, doc.Encoding())
		}
	}

	for _, v := range []string{"1.5", "4.12", "12.5"} {
		doc.SetVersion(v)
		if doc.Version() != v {
			t.Errorf("Expected version '%s', got '%s'", v, doc.Version())
		}
	}

	doc.SetStandalone(1)
	if doc.Standalone() != 1 {
		t.Errorf("Expected standalone 1, got '%d'", doc.Standalone())
	}

	doc.SetBaseURI("localhost/here.xml")
	if doc.URI() != "localhost/here.xml" {
		t.Errorf("Expected URI 'localhost/here.xml', got '%s'", doc.URI())
	}
}

func checkElement(t *testing.T, e *Element, assertName, testCase string) bool {
	if e == nil {
		t.Errorf("%s: Element is nil", testCase)
		return false
	}

	if e.NodeType() != ElementNode {
		t.Errorf("%s: Expected node type 'ElementNode', got '%s'", testCase, e.NodeType())
		return false
	}

	if e.NodeName() != assertName {
		t.Errorf("%s: Expected NodeName '%s', got '%s'", testCase, assertName, e.NodeName())
		return false
	}
	return true
}

func createElementAndCheck(t *testing.T, doc *Document, name, assertName, testCase string) bool {
	node, err := doc.CreateElement(name)
	if err != nil {
		t.Errorf("Failed to create new element '%s': %s", name, err)
		return false
	}
	return checkElement(t, node, assertName, testCase)
}

func withDocument(cb func(*Document)) {
	doc := CreateDocument()
	defer doc.Free()

	cb(doc)
}

func TestDocumentCreateElements(t *testing.T) {
	withDocument(func(d *Document) {
		createElementAndCheck(t, d, "foo", "foo", "Simple Element")
	})

	withDocument(func(d *Document) {
		d.SetEncoding("iso-8859-1")
		createElementAndCheck(t, d, "foo", "foo", "Create element with document with encoding")
	})

	withDocument(func(d *Document) {
		caseName := "Create element with namespace"
		e, err := d.CreateElementNS("http://kungfoo", "foo:bar")
		if err != nil {
			t.Errorf("failed to create namespaced element: %s", err)
			return
		}

		checkElement(t, e, "foo:bar", caseName)

		if e.Prefix() != "foo" {
			t.Errorf("%s: Expected prefix '%s', got '%s'", caseName, "foo", e.Prefix())
		}
		if e.LocalName() != "bar" {
			t.Errorf("%s: Expected local name '%s', got '%s'", caseName, "bar", e.LocalName())
		}
		if e.NamespaceURI() != "http://kungfoo" {
			t.Errorf("%s: Expected namespace uri '%s', got '%s'", caseName, "http://kungfoo", e.NamespaceURI())
		}
	})

	// Bad elements
	withDocument(func(d *Document) {
		badnames := []string{";", "&", "<><", "/", "1A"}
		for _, name := range badnames {
			if _, err := d.CreateElement(name); err == nil {
				t.Errorf("Creation of element name '%s' should fail", name)
			}
		}
	})
}

func TestDocumentCreateText(t *testing.T) {
	withDocument(func(d *Document) {
		node, err := d.CreateTextNode("foo")
		if err != nil {
			t.Errorf("Failed to create text node: %s", err)
			return
		}

		if node.NodeType() != TextNode {
			t.Errorf("Expected NodeType '%s', got '%s'", TextNode, node.NodeType())
			return
		}

		if node.NodeValue() != "foo" {
			t.Errorf("Expeted NodeValue 'foo', got '%s'", node.NodeValue())
			return
		}
	})
}

func TestDocumentCreateComment(t *testing.T) {
	withDocument(func(d *Document) {
		node, err := d.CreateCommentNode("foo")
		if err != nil {
			t.Errorf("Failed to create Comment node: %s", err)
			return
		}

		if node.NodeType() != CommentNode {
			t.Errorf("Expected NodeType '%s', got '%s'", CommentNode, node.NodeType())
			return
		}

		if node.NodeValue() != "foo" {
			t.Errorf("Expeted NodeValue 'foo', got '%s'", node.NodeValue())
			return
		}

		if node.String() != "<!--foo-->" {
			t.Errorf("Expeted String() to return 'foo', got '%s'", node.String())
			return
		}
	})
}

func TestDocumentCreateCDataSection(t *testing.T) {
	withDocument(func(d *Document) {
		node, err := d.CreateCDataSection("foo")
		if err != nil {
			t.Errorf("Failed to create CDataSection node: %s", err)
			return
		}

		if node.NodeType() != CDataSectionNode {
			t.Errorf("Expected NodeType '%s', got '%s'", CDataSectionNode, node.NodeType())
			return
		}

		if node.NodeValue() != "foo" {
			t.Errorf("Expeted NodeValue 'foo', got '%s'", node.NodeValue())
			return
		}

		if node.String() != "<![CDATA[foo]]>" {
			t.Errorf("Expeted String() to return 'foo', got '%s'", node.String())
			return
		}
	})
}

func TestDocumentCreateAttribute(t *testing.T) {
	withDocument(func(d *Document) {
		node, err := d.CreateAttribute("foo", "bar")
		if err != nil {
			t.Errorf("Failed to create Attribute node: %s", err)
			return
		}

		if node.NodeType() != AttributeNode {
			t.Errorf("Expected NodeType '%s', got '%s'", AttributeNode, node.NodeType())
			return
		}

		if node.NodeName() != "foo" {
			t.Errorf("Expeted NodeName 'foo', got '%s'", node.NodeName())
			return
		}

		if node.NodeValue() != "bar" {
			t.Errorf("Expeted NodeValue 'foo', got '%s'", node.NodeValue())
			return
		}

		if node.String() != ` foo="bar"` {
			t.Errorf(`Expeted String() to return ' foo="bar"', got '%s'`, node.String())
			return
		}

		if node.HasChildNodes() {
			t.Errorf("Expected HashChildNodes to return false")
			return
		}

		// Attribute nodes claim to not have any child nodes, but they do?!
		content := node.FirstChild()
		if content == nil {
			t.Errorf("Expected FirstChild to return a node")
			return
		}

		if content.NodeType() != TextNode {
			t.Errorf("Expected content node NodeType '%s', got '%s'", TextNode, content.NodeType())
			return
		}
	})

	// Bad elements
	withDocument(func(d *Document) {
		badnames := []string{";", "&", "<><", "/", "1A"}
		for _, name := range badnames {
			if _, err := d.CreateAttribute(name, "bar"); err == nil {
				t.Errorf("Creation of attribute name '%s' should fail", name)
			}
		}
	})
}

func TestDocumentCreateAttributeNS(t *testing.T) {
	withDocument(func(d *Document) {
		elem, err := d.CreateElement("foo")
		if err != nil {
			t.Errorf("Failed to create Element node: %s", err)
			return
		}
		d.SetDocumentElement(elem)

		attr, err := d.CreateAttribute("attr", "e & f")
		if err != nil {
			t.Errorf("Failed to create Attribute node: %s", err)
			return
		}
		elem.AddChild(attr)

		if elem.String() != `<foo attr="e &amp; f"/>` {
			t.Errorf(`Expected String '<foo attr="e &amp; f"/>', got '%s'`, elem.String())
			return
		}
		elem.RemoveAttribute("attr")

		attr, err = d.CreateAttributeNS("", "attr2", "a & b")
		if err != nil {
			t.Errorf("Failed to create Attribute node: %s", err)
			return
		}
		elem.AddChild(attr)

		if elem.String() != `<foo attr2="a &amp; b"/>` {
			t.Errorf(`Expected String '<foo attr2="a &amp; b"/>', got '%s'`, elem.String())
			return
		}
		elem.RemoveAttribute("attr2")

		attr, err = d.CreateAttributeNS("http://kungfoo", "foo:attr3", "g & h")
		if err != nil {
			t.Errorf("Failed to create Attribute node: %s", err)
			return
		}
		elem.AddChild(attr)

		if elem.String() != `<foo xmlns:foo="http://kungfoo" foo:attr3="g &amp; h"/>` {
			t.Errorf(`Expected String '<foo xmlns:foo="http://kungfoo" foo:attr3="g &amp; h"/>', got '%s'`, elem.String())
			return
		}
	})
}

